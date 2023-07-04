package notify

import (
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (c *Client) msgHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("mqtt msg received; topic: %v\n", msg.Topic())
	for _, v := range c.subs {
		if strings.HasPrefix(msg.Topic(), v.topic) {
			v.cb(msg.Topic(), msg.Payload())
		}
	}
}

func (c *Client) Subscribe(topic string, cb TopicCBFunc) error {
	c.subs = append(c.subs, subCB{
		topic: topic,
		cb:    cb,
	})
	return nil
}
