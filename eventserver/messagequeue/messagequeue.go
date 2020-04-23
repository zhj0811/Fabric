package messagequeue

import (
	"fmt"
	"strings"
	"time"

	"github.com/zhj0811/fabric/common/metadata"

	"github.com/op/go-logging"
	"github.com/streadway/amqp"
)

type mqInfo struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
	address    string
	queueName  string
	isAvaiable bool
}

var (
	mqs    []*mqInfo
	mqsNum int
	logger = logging.MustGetLogger(metadata.LogModule)
)

const (
	RabbitmqReconnectCount = 3
	MaxSendTime            = 3
)

func InitMQ(queueName string, addresses ...string) bool {
	var isAvailableMQ = false
	for _, addr := range addresses {
		mq := mqInfo{
			address:    addr,
			queueName:  queueName,
			isAvaiable: false,
		}
		if mq.connect() {
			isAvailableMQ = true
			logger.Infof("open connect to the rabbitmq address(%s) success!", addr)
		} else {
			logger.Errorf("open connect to the rabbitmq address(%s) failed!", addr)
		}
		mqs = append(mqs, &mq)
	}

	mqsNum = len(mqs)
	logger.Infof("There is %d configured mq.", mqsNum)

	return isAvailableMQ
}

func Close() {
	for _, mq := range mqs {
		if nil != mq.conn {
			mq.conn.Close()
		}
	}
}

func SendMessage(msg interface{}) error {
	var data []byte
	var searchAvaialbe = true

	if ret, ok := msg.([]byte); ok {
		data = ret
	} else if ret, ok := msg.(string); ok {
		data = []byte(ret)
	} else {
		logger.Error("The msg is unexpect type!")
		return fmt.Errorf("Unexpect msg type !")
	}
	if len(data) == 0 {
		logger.Error("The message is empty, nothing to be !")
		return fmt.Errorf("The message is empty!")
	}

	// first use first available mq which was used successfully
	// if mq becomes unavailable, try to use all the addresses to find an available connection
	for i := 0; i < mqsNum; {
		mq := mqs[i]
		if searchAvaialbe && false == mq.isAvaiable {
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
				mq.isAvaiable = false
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
	m.isAvaiable = true

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
	if !m.isAvaiable {
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
