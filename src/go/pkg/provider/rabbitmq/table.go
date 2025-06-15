package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type Table map[string]interface{}

func tableToAMQPTable(table Table) amqp.Table {
	new := amqp.Table{}
	for k, v := range table {
		new[k] = v
	}
	return new
}
