package sdk

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/peersafe/tradetrain/define"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

func GetKafkaNumber() (int, error) {
	topic := viper.GetString("kafka.topic")
	address := viper.GetString("kafka.address")
	targetNum := viper.GetInt("kafka.targetnum")
	if "" == topic || "" == address {
		logger.Error("topic or address for kafka is empty, please set it!")
		return define.KafkaConfigFailed, fmt.Errorf("topic or address is empty")
	}

	config := kafkaConfig()
	if nil == config {
		logger.Error("get config for kafka failed.")
		return define.KafkaConfigFailed, fmt.Errorf("get config for kafka failed")
	}

	broker := sarama.NewBroker(address)
	err := broker.Open(config)
	if err != nil {
		logger.Error("connect to kafka failed: %s", err.Error())
		return define.KafkaConnectionFailed, fmt.Errorf("open connection to kafka failed")
	}
	defer broker.Close()

	request := sarama.MetadataRequest{Topics: []string{topic}}
	response, err := broker.GetMetadata(&request)
	if err != nil {
		logger.Errorf("get metadata from kafka failed: %s", err.Error())
		return define.KafkaConnectionFailed, fmt.Errorf("get metadata from kafka failed")
	}

	brokerNum := len(response.Brokers)
	logger.Debugf("The kafka has %d broker now.", brokerNum)
	if brokerNum != targetNum {
		logger.Errorf("The target broker num is %d, but get %d broker currently!", targetNum, brokerNum)
		logger.Error("---------------------------------------------------------------------------------------")
		logger.Error("Please pay attention to the kafka cluster, restore the environment as soon as possible!")
		logger.Error("---------------------------------------------------------------------------------------")
		return define.KafkaBrokerAbnormal, fmt.Errorf("some broker may be stopped")
	}

	return define.KafkaNormal, nil
}

func kafkaConfig() *sarama.Config {
	kafkaTlsEnabled := viper.GetBool("kafka.tls.enabled")
	privateKeyPath := viper.GetString("kafka.tls.privatekeypath")
	certificatePath := viper.GetString("kafka.tls.certificatepath")
	rootcasPath := viper.GetStringSlice("kafka.tls.rootcaspath")

	config := sarama.NewConfig()

	if kafkaTlsEnabled {
		config.Net.TLS.Enable = true

		keyPair, err := tls.LoadX509KeyPair(certificatePath, privateKeyPath)
		if err != nil {
			logger.Errorf("kafka tls load failed: %s", err.Error())
			return nil
		}
		rootCAs := x509.NewCertPool()

		for _, certificate := range rootcasPath {
			caCert, err := ioutil.ReadFile(certificate)
			if err != nil {
				logger.Error("Unable to load CA cert file.")
				return nil
			}
			if !rootCAs.AppendCertsFromPEM(caCert) {
				logger.Error("Unable to parse the root certificate authority certificates (Kafka.Tls.RootCAs)")
				return nil
			}
		}
		config.Net.TLS.Config = &tls.Config{
			Certificates: []tls.Certificate{keyPair},
			RootCAs:      rootCAs,
			MinVersion:   tls.VersionTLS12,
			MaxVersion:   0, // Latest supported TLS version
		}
		config.Net.TLS.Config.ServerName = viper.GetString("kafka.tls.serverhostoverride")
	}

	return config
}
