package env

type QueueConfig struct {
	Broker        string
	DefaultQueue  string
	ResultBackend string
	AMQP          struct {
		Exchange     string
		ExchangeType string
		BindingKey   string
	}
}
