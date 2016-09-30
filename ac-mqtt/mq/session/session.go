package session

import "github.com/surgemq/message"

type Session struct {
	id     string
	topics map[string]byte

	// 遗愿消息
	Will *message.PublishMessage
}

func New(id string) *Session {
	return &Session{id: id, topics: make(map[string]byte)}
}

func (s *Session) Id() string {
	return s.id
}

func (s *Session) AddTopic(topic string, qos byte) {
	s.topics[topic] = qos
}

func (s *Session) SetWillMessage(qos byte, topic, payload []byte, retain bool) {
	s.Will = message.NewPublishMessage()
	s.Will.SetQoS(qos)
	s.Will.SetTopic(topic)
	s.Will.SetPayload(payload)
	s.Will.SetRetain(retain)
}
