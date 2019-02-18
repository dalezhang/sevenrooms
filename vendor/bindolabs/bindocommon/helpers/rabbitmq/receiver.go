package rabbitmq

// 消费者的抽象
type Receiver interface {
	QueueName() string     // 队列名
	RouterKey() string     // 队列路由
	OnError(error)         // 队列错误处理
	OnReceive([]byte) bool // 队列处理的业务函数
}
