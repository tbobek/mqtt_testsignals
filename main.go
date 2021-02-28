package main

import (
	"fmt"
	"math/rand"
	"mqtt_testsignals/mqtthandler"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// HelloServer is an example API endpoint
func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func makenoise(noise float64, value float64, number int) []float64 {
	if number < 1 {
		return []float64{}
	}
	var data1 = make([]float64, number)
	for i := 0; i < number; i++ {
		data1[i] = value + 2*(rand.Float64()-0.5)*noise
	}
	return data1
}

func makeslope(val1 float64, val2 float64, noise float64, number int) []float64 {
	if number < 2 { // to few numbers
		return []float64{}
	}
	var data1 = make([]float64, number)
	span := val2 - val1
	dy := span / float64(number-1)
	for i := 0; i < number; i++ {
		data1[i] = val1 + 2.0*(rand.Float64()-0.5)*noise + dy*float64(i)
	}
	return data1
}

func InitSequence(w http.ResponseWriter, r *http.Request) {
	// settings
	noise := 0.5
	value1 := 1.0
	value2 := 10.0
	// make initial data
	data := makenoise(noise, value1, 10)
	data = append(data, makeslope(value1, value2, noise, 5)...)
	data = append(data, makenoise(noise, value2, 10)...)
	data = append(data, makeslope(value2, value1, noise, 5)...)
	data = append(data, makenoise(noise, value1, 10)...)

	// publish
	mqtthandler.PublishSequence(client, "test/data", data, 1)
}

var client mqtt.Client

func main() {
	var broker = "localhost"
	var port = 8883
	client = mqtthandler.Makeclient(broker, port)
	//go mqtthandler.Publish(client, "topic/test", 10)
	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/seq", InitSequence)
	http.ListenAndServeTLS(":9442", "srv.crt", "srv.key", nil)
	fmt.Println("vim-go")
}
