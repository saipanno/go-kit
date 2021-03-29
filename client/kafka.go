package client

import (
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/saipanno/go-kit/logger"
	kafkacg "github.com/wvanbergen/kafka/consumergroup"
)

// CreateKafkaProducer ...
func CreateKafkaProducer(conf *KafkaConfig) (producer sarama.AsyncProducer, err error) {

	logger.Infof("create kafka connect %s",
		strings.Join(conf.Host, ","))

	kc := sarama.NewConfig()
	kc.Producer.Compression = sarama.CompressionNone
	kc.Producer.Timeout = 10 * time.Second
	kc.Producer.RequiredAcks = sarama.WaitForAll
	kc.Producer.Partitioner = sarama.NewHashPartitioner
	kc.Producer.Return.Successes = false
	kc.Producer.Return.Errors = false
	if conf.UseV010Version {
		kc.Version = sarama.V0_10_0_0
	}

	var client sarama.Client
	client, err = sarama.NewClient(conf.Host, kc)
	if err != nil {
		logger.Errorf("create kafka client failed, message is %s",
			err.Error())
		return
	}

	producer, err = sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		logger.Errorf("create kafka async producer failed, message is %s",
			err.Error())
	}

	return
}

// CreateKafkaCGConn ...
func CreateKafkaCGConn(conf *KafkaConfig) (kafkaCGroup *kafkacg.ConsumerGroup, err error) {

	config := kafkacg.NewConfig()
	config.Offsets.Initial = sarama.OffsetNewest
	config.ClientID = conf.Topic

	kafkaCGroup, err = kafkacg.JoinConsumerGroup(
		conf.Topic,
		[]string{conf.Topic},
		conf.Zookeeper,
		config)
	if err != nil {
		logger.Errorf("create kafka %s ConsumerGroup failed, message is %s",
			conf.Topic, err.Error())
	}

	return
}
