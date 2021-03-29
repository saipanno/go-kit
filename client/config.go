package client

import mqtt "github.com/eclipse/paho.mqtt.golang"

// DBConfig ...
type DBConfig struct {
	URI      string `json:"uri,omitempty"`
	Database string `json:"database,omitempty"`
	// MaxIdle   int    `json:"max_idle"`
	// MaxActive int    `json:"max_active"`
	// Wait      bool   `json:"wait"`
}

// MQTTConfig ...
type MQTTConfig struct {
	Name                  string                `json:"name,omitempty"`
	Broker                []string              `json:"broker,omitempty"`
	QOS                   int                   `json:"qos,omitempty"`
	ClientID              string                `json:"client_id,omitempty"`
	Username              string                `json:"username,omitempty"`
	Password              string                `json:"password,omitempty"`
	CleanSession          bool                  `json:"clean_session,omitempty"`
	OnConnectHandler      mqtt.OnConnectHandler `json:"-"`
	DefaultPublishHandler mqtt.MessageHandler   `json:"-"`
}

// KafkaConfig ...
type KafkaConfig struct {
	Topic          string   `json:"topic,omitempty"`
	Host           []string `json:"host,omitempty"`
	Zookeeper      []string `json:"zookeeper,omitempty"`
	UseV010Version bool     `json:"use_v010_version,omitempty"`
}
