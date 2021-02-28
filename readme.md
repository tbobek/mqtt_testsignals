# Readme

Simple https server in golang, with a MQTTS client that can 
publish a test sequence of numbers triggered by the https server. 

## Release notes

### v0.1 (2021-02-28)
- first version with basic code

## Certificate generation

With help of this tool [https://smallstep.com/hello-mtls/doc/server/go](https://smallstep.com/hello-mtls/doc/server/go) the generation of root certificates and client certificates is very easy. I required the following process to generate certificates for https (output: `root.crt`, `srv.crt`, `srv.key`) and mqtts (output `mqtt_root.crt`, `mqtt.crt`, `mqtt.key`): 

1. `step ca init` to initialize the certificate authority
2. `step-ca ~/.step/config/ca.json` runs a local server as CA authority
3. `step ca root root.crt` to get the root certificate
4. `step ca certificate localhost srv.crt srv.key` to get the server key and certificate


## Usage

Start go program

`go run main.go`

Start subscription on the Mosquitto broker

`mosquitto_sub --cafile mqtt_root.crt -t test/data`

Issue the command via https server

`curl --cacert root.crt https://localhost:9442/seq`

## Hints that I required during development

### How to know if cert is in PEM mode

If you see text output with this command, it is PEM (link): 

`openssl x509 -in cert.crt -text`

### MQTT with TLS

[Code example](https://github.com/eclipse/paho.mqtt.golang/blob/master/cmd/ssl/main.go)

## Todo

- [ ] allow repeated sending of test data
- [ ] multiple test data in own threads
- [ ] Make watchdog routines
- [ ] Configuration via https API
