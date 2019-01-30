package config

// GetRmqConnectionConfig returns rabbitmq cluster host
func (cfg Config) GetRmqConnectionConfig() ConnectionType {
	return cfg.rmqConnection
}

// GetRmqQueueConfig returns rabbitmq cluster host
func (cfg Config) GetRmqQueueConfig(queue string) QueueType {
	return cfg.rmqQueueMap[queue]
}

// Run returns config
func Run() Config {
	return Config{
		rmqConnection: ConnectionType{
			Username: "jsc",
			Password: "qweasdzxc123",
			Host:     "localhost",
			Port:     5673,
			Vhost:    "/",
			Wait:     3,
		},
		rmqQueueMap: map[string]QueueType{
			"cbs_queue_1": QueueType{
				QueueName: "cbs_queue_1",
				Consumer:  "cbs_queue_1",
				AutoAck:   false,
				Exclusive: false,
				NoLocal:   false,
				NoWait:    false,
			},
			"cbs_queue_2": QueueType{
				QueueName: "cbs_queue_2",
				Consumer:  "cbs_queue_2",
				AutoAck:   false,
				Exclusive: false,
				NoLocal:   false,
				NoWait:    false,
			},
		},
	}
}
