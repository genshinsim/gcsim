package notify

import (
	"errors"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	mqttClient mqtt.Client
	subs       []subCB
}

type TopicCBFunc func(topic string, payload []byte)

type subCB struct {
	topic string
	cb    TopicCBFunc
}

func New(id string) (*Client, error) {
	c := &Client{}

	opts, err := c.mqttOpts(id)
	if err != nil {
		return nil, err
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("error connecting to mosquitto: ", token.Error())
		return nil, token.Error()
	}

	c.mqttClient = client

	return c, nil
}

func (c *Client) ActivateListener() error {
	if token := c.mqttClient.Subscribe("#", 1, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *Client) mqttOpts(id string) (*mqtt.ClientOptions, error) {
	host := os.Getenv("MQTT_URL")
	user := os.Getenv("MQTT_USERNAME")
	pass := os.Getenv("MQTT_PASSWORD")
	if host == "" {
		return nil, errors.New("env MQTT_HOST not set")
	}
	if user == "" {
		return nil, errors.New("env MQTT_USER not set")
	}
	if pass == "" {
		return nil, errors.New("env MQTT_PASS not set")
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(host)
	opts.SetUsername(user)
	opts.SetPassword(pass)
	opts.SetClientID(id)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(c.msgHandler)

	return opts, nil
}
