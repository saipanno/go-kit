package client

import (
	"errors"
	"fmt"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/saipanno/go-kit/logger"
)

// CreateMQTTConn ...
func CreateMQTTConn(conf *MQTTConfig, IDPrefix ...string) (client mqtt.Client, err error) {

	ID := conf.ClientID
	if len(IDPrefix) > 0 {
		ID = fmt.Sprintf("%s_%s", IDPrefix[0], ID)
	}

	logger.Infof("create mqtt connect %s with id %s",
		strings.Join(conf.Broker, ","), ID)

	if len(ID) == 0 {
		err = errors.New("mqtt client_id is empty")
		return
	}

	opts := mqtt.NewClientOptions()
	opts.SetClientID(ID).SetPingTimeout(time.Second).SetWriteTimeout(time.Second)

	if !conf.CleanSession {
		opts.SetCleanSession(false)
	}

	if len(conf.Password) > 0 {
		opts.SetUsername(conf.Username).SetPassword(conf.Password)
		if len(conf.Username) == 0 {
			opts.SetUsername(conf.ClientID)
		}
	}

	for _, item := range conf.Broker {
		opts.AddBroker(fmt.Sprintf("tcp://%s", item))
	}

	if conf.OnConnectHandler != nil {
		opts.OnConnect = conf.OnConnectHandler
	}

	if conf.DefaultPublishHandler != nil {
		opts.DefaultPublishHandler = conf.DefaultPublishHandler
	}

	opts.OnConnectionLost = func(_ mqtt.Client, err error) {
		logger.Errorf("mqtt connect lost, message is %s", err.Error())
	}

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		err = token.Error()
	}

	return
}
