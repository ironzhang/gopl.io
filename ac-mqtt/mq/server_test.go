package mq_test

import (
	"ac-common-go/glog"
	"ac-mqtt/mq"
	"fmt"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func StartTestServer() *mq.Server {
	s := mq.NewServer("tcp", "localhost:7000")
	go func() {
		if err := s.ListenAndServe(); err != nil {
			glog.Fatal(err)
		}
	}()
	time.Sleep(50 * time.Millisecond)
	return s
}

func RunTestClient() {
	opts := mqtt.NewClientOptions()
	opts.KeepAlive = 3 * time.Second
	opts.AddBroker("tcp://localhost:7000")

	client := mqtt.NewClient(opts)

	fmt.Println("connect ...")
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("connect: %v\n", token.Error())
		return
	}

	fmt.Println("0 Qos publish ...")
	t0 := client.Publish("topic", 0, false, "payload")

	fmt.Println("1 Qos publish ...")
	t1 := client.Publish("topic", 1, false, "payload")

	fmt.Println("2 Qos publish ...")
	t2 := client.Publish("topic", 2, false, "payload")

	if t0.Wait() {
		fmt.Printf("0 Qos publish: %v\n", t0.Error())
	}

	if t1.Wait() {
		fmt.Printf("1 Qos publish: %v\n", t1.Error())
	}

	if t2.Wait() {
		fmt.Printf("2 Qos publish: %v\n", t2.Error())
	}

	var token mqtt.Token

	fmt.Println("0 Qos subscribe ...")
	token = client.Subscribe("topic", 0, func(c mqtt.Client, m mqtt.Message) {
		fmt.Printf("subscribe: %s\n", m.Topic())
	})
	if token.Wait() {
		fmt.Printf("0 Qos subscribe: %v\n", token.Error())
	}

	fmt.Println("unsubscribe ...")
	token = client.Unsubscribe("topic")
	if token.Wait() {
		fmt.Printf("unsubscribe: %v\n", token.Error())
	}

	fmt.Println("Disconnect ...")
	client.Disconnect(200)
	fmt.Println("Disconnected")
}

func TestServer(t *testing.T) {
	defer glog.Flush()
	go func() {
		tick := time.NewTicker(time.Second)
		for range tick.C {
			glog.Flush()
		}
	}()

	s := StartTestServer()

	RunTestClient()
	//RunTestClient()

	time.Sleep(200 * time.Millisecond)
	s.Close()
	time.Sleep(500 * time.Millisecond)
}
