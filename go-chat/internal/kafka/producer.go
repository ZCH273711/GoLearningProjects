package kafka

import (
	"examples/go-chat/pkg/global/log"
	"github.com/Shopify/sarama"
	"strings"
)

var producer sarama.AsyncProducer
var topic string = "default_message"

func InitProducer(topicInput, hosts string) {
	topic = topicInput
	config := sarama.NewConfig()
	config.Producer.Compression = sarama.CompressionGZIP
	// 连接到kafka集群中的一个服务器，即Broker，并从Broker获取元数据
	// hosts为这些Broker的地址
	client, err := sarama.NewClient(strings.Split(hosts, ","), config)
	if err != nil {
		log.Logger.Error("init kafka client error", log.Any("init kafka client error", err.Error()))
	}

	producer, err = sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		log.Logger.Error("init kafka async client error", log.Any("init kafka async client error", err.Error()))
	}
}

func Send(data []byte) {
	be := sarama.ByteEncoder(data)
	producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: be}
}
