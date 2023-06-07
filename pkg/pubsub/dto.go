package pubsub

import "cloud.google.com/go/pubsub"

type PubSubMessage struct {
	Message      *pubsub.Message
	Subscription string `json:"subscription"`
}

type PubSubPublishDTO struct {
	TopicID     string
	Data        []byte
	OrderingKey *string
}
