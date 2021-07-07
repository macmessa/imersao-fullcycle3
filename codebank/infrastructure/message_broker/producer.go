package message_broker

import kafka "github.com/confluentinc/confluent-kafka-go/kafka"

type MessageProducer struct {
	Producer *kafka.Producer
}

func NewMessageProducer() MessageProducer {
	return MessageProducer{}
}

func (m *MessageProducer) SetupProducer(bootstrapServer string) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
	}
	m.Producer, _ = kafka.NewProducer(configMap)
}

func (m *MessageProducer) Publish(msg string, topic string) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg),
	}
	err := m.Producer.Produce(message, nil)

	if err != nil {
		return err
	}

	return nil
}
