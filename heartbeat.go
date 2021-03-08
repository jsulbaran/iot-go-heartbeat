package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tkanos/gonfig"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/url"
	"time"
)

var (
	err           error
	configuration = Configuration{}
	mqttClient    mqtt.Client
	lastEvent     = ""
)

func exitOnError(e error) {
	if err != nil {
		fmt.Print(err)
		log.Fatal(err)
	}
}
func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}
func connect(clientId string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	exitOnError(token.Error())
	return client
}
func createConnection(configuration Configuration) mqtt.Client {
	uri, err := url.Parse(configuration.ServerUrl)
	exitOnError(err)
	return connect("pub", uri)
}
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}
func main() {
	err := gonfig.GetConf("config.json", &configuration)
	exitOnError(err)

	log.SetOutput(&lumberjack.Logger{
		Filename:   "sonoff.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

	mqttClient = createConnection(configuration)
	doEvery(time.Duration(configuration.HeartbeatFrequencySeconds)*time.Second, sendHeartbeat)
	mqttClient.Disconnect(2000)
}
func sendIot(iotEvent IotEvent, topic string, client mqtt.Client) {
	exitOnError(err)
	jsonValue, err := json.Marshal(iotEvent)
	exitOnError(err)

	if lastEvent == string(jsonValue) {
		return
	}
	lastEvent = string(jsonValue)
	if client.IsConnected() {

		token := client.Publish(topic, 2, true, string(jsonValue))
		token.Wait()
		//client.Disconnect(2000)
	} else {
		fmt.Println("cliente no conectado")
		mqttClient = createConnection(configuration)
	}
}
func sendHeartbeat(t time.Time) {
	timestampStr := time.Now().Format("2006-01-02T15:04:05")
	fmt.Println("Sending heartbeat... : " + timestampStr)
	iotEvent := IotEvent{DeviceId: configuration.DeviceId, Timestamp: timestampStr, EventType: "HEARTBEAT"}
	sendIot(iotEvent, configuration.HeartbeatTopicName, mqttClient)
	fmt.Println("finishing waiting...")
}
