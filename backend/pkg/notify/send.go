package notify

func (c *Client) Notify(topic string, msg any) error {
	token := c.mqttClient.Publish(topic, 0, false, msg)
	<-token.Done()
	return token.Error()
}
