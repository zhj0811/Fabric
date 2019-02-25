package messagequeue

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"

	"github.com/peersafe/factoring/common/metadata"

	"github.com/op/go-logging"
	"github.com/streadway/amqp"
)

type mqInfo struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
	address    string
	queueName  string
	isAvailable bool
}

// TODO: add system mq queue
type mqBaseInfo struct {
	isEnable bool
	addresses []string
	queues []string
	queueForChannels map[string][]string
	channelForQueue map[string]string
}

var (
	mqInfos    = make(map[string][]*mqInfo)
	baseInfo mqBaseInfo
	//mqsNum int
	logger = logging.MustGetLogger(metadata.LogModule)
)

const (
	RabbitmqReconnectCount = 3
	MaxSendTime            = 3
)

func InitMQBaseInfo() bool {
	baseInfo.isEnable = viper.GetBool("mq.mqEnable")
	if !baseInfo.isEnable {
		logger.Infof("mq is disable")
		return true
	}

	if baseInfo.addresses = viper.GetStringSlice("mq.mqAddress"); len(baseInfo.addresses) == 0 {
		logger.Errorf("mq address is null")
		return false
	}
	baseInfo.queueForChannels = make(map[string][]string)
	if baseInfo.queueForChannels =  viper.GetStringMapStringSlice("mq.queues"); len(baseInfo.queueForChannels) == 0 {
		logger.Errorf("mq queue is null")
		return false
	}
	baseInfo.channelForQueue = make(map[string]string)
	for key, values := range baseInfo.queueForChannels {
		logger.Debugf("add %s to mq queues", key)
		baseInfo.queues = append(baseInfo.queues, key)
		for _, value := range values {
			logger.Debugf("the data from channel %s will be sent to queue %s", value, key)
			baseInfo.channelForQueue[value] = key
		}
	}
	logger.Infof("init mq base info successfully")

	return baseInfo.InitMQ()
}

func (baseInfo mqBaseInfo) GetQueueByChannel(channel string) string {
	return baseInfo.channelForQueue[channel]
}

func (baseInfo mqBaseInfo)InitMQ() bool {
	for _, queueName := range baseInfo.queues {
		var isAvailableMQ = false
		logger.Debugf("try to connect to rabbitmq for channel %s", queueName)
		for _, addr := range baseInfo.addresses {
			mq := mqInfo{
				address:    addr,
				queueName:  queueName,
				isAvailable: false,
			}
			if mq.connect() {
				isAvailableMQ = true
				logger.Infof("open connect to the rabbitmq address(%s) for channel(%s) success!",
					addr, queueName)
			} else {
				logger.Errorf("open connect to the rabbitmq address(%s) for channel(%s) failed!",
					addr, queueName)
			}
			mqInfos[queueName] = append(mqInfos[queueName], &mq)
		}
		if !isAvailableMQ {
			logger.Errorf("can connect to the all rabbitmq for queue %s", queueName)
			return false
		}

		logger.Infof("queue %s has %d available rabbitmq connections", queueName, len(mqInfos[queueName]))
	}

	logger.Infof("init mq successfully")
	return true
}

func Close() {
	for k, mqs := range mqInfos {
		logger.Debugf("try to close queue %s's mq", k)
		for _, mq := range mqs {
			if nil != mq.conn {
				logger.Debugf("try to close connection to mq %s", mq.address)
				if err := mq.conn.Close(); err != nil {
					logger.Errorf("close connection to mq %s failed: %s", mq.address, err.Error())
				} else {
					logger.Debugf("close connection to mq %s successfully", mq.address)
				}
			} else {
				logger.Debugf("the connection to mq %s has not been initialized", mq.address)
			}
		}
	}
	logger.Infof("close all mq connections successfully")
}

func SendMessage(channelName string, msg interface{}) error {
	var data []byte
	var searchAvaialbe = true

	if ret, ok := msg.([]byte); ok {
		data = ret
	} else if ret, ok := msg.(string); ok {
		data = []byte(ret)
	} else {
		logger.Error("The msg is unexpected type!")
		return fmt.Errorf("Unexpect msg type !")
	}
	if len(data) == 0 {
		logger.Error("The message is empty, nothing to be sent!")
		return fmt.Errorf("The message is empty.")
	}

	queueName := baseInfo.GetQueueByChannel(channelName)
	mqs := mqInfos[queueName]
	mqsNum := len(mqs)

	// first use first available mq which was used successfully
	// if mq becomes unavailable, try to use all the addresses to find an available connection
	for i := 0; i < mqsNum; {
		mq := mqs[i]
		if searchAvaialbe && false == mq.isAvailable {
			logger.Infof("mq(%s) is invalid, looking for the next one.", mq.address)
			i++
			continue
		}

		tryTime := 0
		for {
			if err := mq.send(data); err == nil {
				return nil
			} else {
				logger.Errorf("send data to mq %s failed: %s", mq.address, err.Error())
				mq.isAvailable = false
				tryTime++
				if tryTime >= MaxSendTime {
					logger.Errorf("send message to mq %s has reached max time!", mq.address)
					searchAvaialbe = false
					break
				}
				// Exception (504) Reason: "channel/connection is not open"
				// If the connection is unavailable, try to reconnect to it.
				if strings.Contains(err.Error(), "channel/connection is not open") && mq.reConnect() {
					logger.Warningf("reconnect to rabbitmq %s successfully and try to resend it!", mq.address)
					continue
				} else {
					if searchAvaialbe {
						logger.Warningf("mq(%s) is unavailable, try to find available mq from all.", mq.address)
						searchAvaialbe = false
						i = 0
					} else {
						logger.Warningf("mq(%s) is unavailable, try to use the next one!", mq.address)
						i++
					}
					break
				}
			}
		}
	}

	return fmt.Errorf("Send message to all mq failed!")
}

func (m *mqInfo) connect() bool {
	conn, err := amqp.Dial(m.address)
	if err != nil {
		logger.Errorf("Failed to connect to RabbitMQ: %s", err.Error())
		return false
	}

	channel, err := conn.Channel()
	if err != nil {
		logger.Errorf("Failed to open a channel: %s", err.Error())
		conn.Close()
		return false
	}

	queue, err := channel.QueueDeclare(
		m.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		logger.Errorf("Failed to declare a queue: %s", err.Error())
		conn.Close()
		return false
	}

	//release the old connection.
	if nil != m.conn {
		m.conn.Close()
	}

	m.conn = conn
	m.channel = channel
	m.queue = queue
	m.isAvailable = true

	return true
}

func (m *mqInfo) reConnect() bool {
	for sendTime := 0; sendTime < RabbitmqReconnectCount; sendTime++ {
		if m.connect() {
			logger.Infof("Reconnect to rabbitmq %s successfully!", m.address)
			return true
		}
		time.Sleep(time.Second)
	}

	logger.Errorf("Reconnect to rabbitmq %s has reached the max time.", m.address)
	return false
}

//if the mq is marked as unavailable, first to establish connection.
func (m *mqInfo) send(sendData []byte) error {
	if !m.isAvailable {
		logger.Warningf("The mq %s is not available, but still try to send data.", m.address)
		if m.connect() {
			logger.Infof("reconnect to mq %s successfully, try to send data.", m.address)
		} else {
			logger.Errorf("reconnect to mq %s failed, not to send data.", m.address)
			return fmt.Errorf("channel/connection is not open")
		}
	}

	if nil == m.channel {
		logger.Error("mq %s doesn't have channel.", m.address)
		return fmt.Errorf("channel/connection is not open")
	}

	err := m.channel.Publish(
		"",           // exchange
		m.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        sendData,
		})

	return err
}
