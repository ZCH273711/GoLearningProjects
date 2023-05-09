package kafka

import (
	"examples/go-chat/pkg/global/log"
	"github.com/Shopify/sarama"
	"strings"
)

var consumer sarama.Consumer

type ConsumerCallBack func(data []byte)

func InitConsumer(hosts string) {
	config := sarama.NewConfig()
	client, err := sarama.NewClient(strings.Split(hosts, ","), config)
	if err != nil {
		log.Logger.Error("init kafka consumer client error", log.Any("init kafka consumer client error", err.Error()))
	}

	consumer, err = sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Logger.Error("init kafka consumer client error", log.Any("init kafka consumer client error", err.Error()))

	}
}

func ConsumerMsg(callBack ConsumerCallBack) {
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if nil != err {
		log.Logger.Error("iConsumePartition error", log.Any("ConsumePartition error", err.Error()))
		return
	}

	defer partitionConsumer.Close()
	for {
		msg := <-partitionConsumer.Messages()
		if nil != callBack {
			callBack(msg.Value)
		}
	}
}

func CloseConsumer() {
	if consumer != nil {
		consumer.Close()
	}
}
