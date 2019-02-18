package rabbitmq

import (
	"bindolabs/golib/logger"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

var publishChan *amqp.Channel // 生产者的管道

type Channel struct {
	wg           sync.WaitGroup
	exchangeName string        // 交换机名称
	exchangeType string        // 交换机类型
	args         amqp.Table    // 交换机其他设置
	receiverChan *amqp.Channel // 消费者的管道
	receivers    []Receiver    // 消费者
}

// 构建交换机
//
// c: 已经连接上的RabbitMQ connection
// exchangeName: 交换机名称
// exchangeType: 交换机类型
// args: 交换机的其他配置参数
func NewChannel(c *Connection, exchangeName, exchangeType string, args amqp.Table) (*Channel, error) {
	var err error
	channel := &Channel{
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		args:         args,
	}
	channel.receiverChan, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}
	err = channel.prepareExchange()
	if err != nil {
		return nil, err
	}
	return channel, nil
}

// 构建Exchange
func (c *Channel) prepareExchange() error {
	return c.receiverChan.ExchangeDeclare(
		c.exchangeName, // exchange
		c.exchangeType, // type
		true,           // 持久化
		false,          // 自动删除
		false,          // internal
		false,          // noWait
		c.args,         // args
	)
}

// 注册消费者服务
func (c *Channel) RegisterReceiver(r Receiver) {
	c.receivers = append(c.receivers, r)
}

// 注册多个消费者服务
func (c *Channel) RegisterReceivers(rs []Receiver) {
	for i := range rs {
		c.RegisterReceiver(rs[i])
	}
}

// 开启队列
// 建议使用时考虑网络等异常中断的重试
func (c *Channel) RunReceivers() {
	for _, receiver := range c.receivers {
		c.wg.Add(1)
		go c.listen(receiver) // 每个接收者单独启动一个goroutine用来初始化queue并接收消息
	}
	c.wg.Wait()

	logger.Errorf("所有处理queue的任务都意外退出了")
}

// Listen 监听指定路由发来的消息
// 这里需要针对每一个接收者启动一个goroutine来执行listen
// 该方法负责从每一个接收者监听的队列中获取数据，并负责重试
func (c *Channel) listen(receiver Receiver) {
	defer func() {
		c.wg.Done()
		if err := recover(); err != nil {
			stack := debug.Stack()
			logger.Errorf("receivers panic: %v\n%s\n", err, stack)
		}
	}()

	// 这里获取每个接收者需要监听的队列和路由
	queueName := receiver.QueueName()
	routerKey := receiver.RouterKey()

	// 申明Queue
	_, err := c.receiverChan.QueueDeclare(
		queueName, // name
		true,      // durable (持久化)
		false,     // delete when usused (自动删除)
		false,     // exclusive(排他性队列)
		false,     // no-wait
		nil,       // arguments
	)
	if nil != err {
		// 当队列初始化失败的时候，需要告诉这个接收者相应的错误
		receiver.OnError(fmt.Errorf("初始化队列 %s 失败: %s", queueName, err.Error()))
	}

	// 将Queue绑定到Exchange上去
	err = c.receiverChan.QueueBind(
		queueName,      // queue name
		routerKey,      // routing key
		c.exchangeName, // exchange
		false,          // no-wait
		nil,
	)
	if nil != err {
		receiver.OnError(fmt.Errorf("绑定队列 [%s - %s] 到交换机失败: %s", queueName, routerKey, err.Error()))
	}

	c.receiverChan.Qos(1, 0, true) // 确保rabbitmq会一个一个发消息
	msgs, err := c.receiverChan.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if nil != err {
		receiver.OnError(fmt.Errorf("获取队列 %s 的消费通道失败: %s", queueName, err.Error()))
	}

	var retryTimes = 0 // 重试次数记录器，最大次数为10

	// 使用callback消费数据
	for msg := range msgs {
		// 当接收者消息处理失败的时候，
		// 比如网络问题导致的数据库连接失败，redis连接失败等等这种
		// 通过重试可以成功的操作，那么这个时候是需要重试的
		// 直到数据处理成功后再返回，然后才会回复rabbitmq ack
		for !receiver.OnReceive(msg.Body) && retryTimes < 10 {
			retryTimes += 1
			logger.Warnf("receiver 数据处理失败，将要重试")
			time.Sleep(1 * time.Second)
		}

		// TODO: 这里如果处理失败10次，应该记录message，而不是直接ack

		// 确认收到本条消息, multiple必须为false
		msg.Ack(false)
	}
}
