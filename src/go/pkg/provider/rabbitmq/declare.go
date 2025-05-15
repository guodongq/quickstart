package rabbitmq

func declareQueue(conn *Connection, options *QueueOptions) error {
	if !options.Declare {
		return nil
	}
	if options.Passive {
		queue, err := conn.ch.QueueDeclarePassive(
			options.Name,
			options.Durable,
			options.AutoDelete,
			options.Exclusive,
			options.NoWait,
			tableToAMQPTable(options.Args),
		)
		if err != nil {
			return err
		}
		if len(options.Name) == 0 {
			options.TemporaryName = queue.Name
		}
		return err
	}

	queue, err := conn.ch.QueueDeclare(
		options.Name,
		options.Durable,
		options.AutoDelete,
		options.Exclusive,
		options.NoWait,
		tableToAMQPTable(options.Args),
	)
	if err != nil {
		return err
	}
	if len(options.Name) == 0 {
		options.TemporaryName = queue.Name
	}
	return nil
}

func declareExchange(conn *Connection, options ExchangeOptions) error {
	if !options.Declare {
		return nil
	}
	if options.Passive {
		err := conn.ch.ExchangeDeclarePassive(
			options.Name,
			options.Kind,
			options.Durable,
			options.AutoDelete,
			options.Internal,
			options.NoWait,
			tableToAMQPTable(options.Args),
		)
		return err
	}

	err := conn.ch.ExchangeDeclare(
		options.Name,
		options.Kind,
		options.Durable,
		options.AutoDelete,
		options.Internal,
		options.NoWait,
		tableToAMQPTable(options.Args),
	)
	return err
}

func declareBindings(conn *Connection, options ConsumerOptions) error {
	for _, exchangeOption := range options.ExchangeOptions {
		for _, binding := range exchangeOption.Bindings {
			if !binding.Declare {
				continue
			}
			err := conn.ch.QueueBind(
				options.QueueOptions.Name,
				binding.RoutingKey,
				exchangeOption.Name,
				binding.NoWait,
				tableToAMQPTable(binding.Args),
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
