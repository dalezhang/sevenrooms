package rabbitmq

import (
	"bindolabs/bindocommon/env"
	"fmt"

	"github.com/streadway/amqp"
)

// 交换机类型: 延迟队列
const DELAY_TYPE = "x-delayed-message"

// RabbitMQ 连接对象
// 它主要负责建立连接，使用连接建立交换机
type Connection struct {
	conn *amqp.Connection
}

// 创建 RabbitMQ 连接对象
func New() *Connection {
	conn, err := amqp.Dial(env.Env.RabbitMQ.DNS)
	if err != nil {
		panic(fmt.Sprintf("RabbitMQ连接失败: %s", err.Error()))
	}
	return &Connection{
		conn: conn,
	}
}

// 使用连接创建 Topic 类型的交换机
//
// exchangeName: 交换机名称
// args: 交换机的其他配置参数
func (c *Connection) NewTopic(exchangeName string, args amqp.Table) (*Channel, error) {
	return NewChannel(c, exchangeName, amqp.ExchangeTopic, nil)
}

// 使用连接创建 Delay 类型的交换机
//
// exchangeName: 交换机名称
// DelayedType: RabbitMQ的官方交换机类型，交换机会在这个交换类型上做延迟
// args: 交换机的其他配置参数
func (c *Connection) NewDelay(exchangeName, DelayedType string, args amqp.Table) (*Channel, error) {
	return NewChannel(c, exchangeName, DELAY_TYPE, amqp.Table{"x-delayed-type": DelayedType})
}

// 生产消息
//
// exchangeName: 交换机名称
// routerKey: 队列路由
// Publishing: 发送数据
func (c *Connection) Publish(exchangeName, routerKey string, p amqp.Publishing) error {
	// 单例，所有的publish都是由 publishChan 负责
	var err error
	if publishChan == nil {
		publishChan, err = c.conn.Channel()
		if err != nil {
			return err
		}
	}

	if len(p.ContentType) < 1 {
		p.ContentType = "text/plain"
	}
	if p.DeliveryMode != 0 && p.DeliveryMode != 1 && p.DeliveryMode != 2 {
		p.DeliveryMode = 2
	}
	return publishChan.Publish(exchangeName, routerKey, false, false, p)
}
