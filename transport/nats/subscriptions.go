package nats

import "encoding/json"

type subData struct {
	Data interface{}
}

func (p *Publisher) publish(subject string, data interface{}) {
	encoded, err := json.Marshal(subData{
		Data: data,
	})
	if err != nil {
		p.Log.Printf("error when encoding nats message for publishing to subject %s: %s\n", subject, err)
		return
	}

	if err := p.Conn.Publish(subject, encoded); err != nil {
		p.Log.Printf("error when publishing nats message to subject %s: %s\n", subject, err)
		return
	}
}
