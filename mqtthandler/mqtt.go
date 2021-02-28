package mqtthandler

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Makeclient(broker string, port int) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	if port == 8883 {
		tlsConfig := NewTlsConfig()
		opts.SetTLSConfig(tlsConfig)
	}
	//opts.SetUsername("emqx")
	//opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}

func NewTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("mqtt_root.crt")
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)
	// Import client certificate/key pair
	//clientKeyPair, err := tls.LoadX509KeyPair("mqtt.crt", "mqtt.key")
	//if err != nil {
	//	panic(err)
	//}
	return &tls.Config{
		RootCAs: certpool,
		//ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		//Certificates:       []tls.Certificate{clientKeyPair},
	}
}

// Publish does some publishing
func Publish(client mqtt.Client, topic string, num int) {
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish(topic, 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

// PublishSequence outputs a number sequence on a specific MQTT topic
func PublishSequence(client mqtt.Client, topic string, data []float64, sleeptime_sec float64) {
	for _, d := range data {
		token := client.Publish(topic, 0, false, fmt.Sprintf("%f", d))
		token.Wait()
		time.Sleep(time.Second * time.Duration(sleeptime_sec))
	}
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}
