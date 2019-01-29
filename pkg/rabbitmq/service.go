package rabbitmq

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/streadway/amqp"
)

func reConnect(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("reConnect cleanup..\n")
			close(cleanup["reconnect"])
			return
		case <-rmq.startConnection:
			if rmq.reconn.Load().(bool) {
				continue
			} else {
				rmq.createConnect()
			}
		}
	}
}

func reChannel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("reChannel cleanup..\n")
			close(cleanup["rechannel"])
			return
		case <-rmq.startChannel:
			if rmq.rechan.Load().(bool) {
				continue
			} else {
				rmq.createChannel()
			}
		}
	}
}

func catchAmqpEvent(ctx context.Context) {
	// Only rmq.startConnection <- struct{}{} can be called on receiving
	// event from the library.
	// Try to use other channel would cause dead lock.
	// By design there's only single entry for the chain reaction.
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("catchChannelEvent cleanup..\n")
			close(cleanup["catchevent"])
			return
		case err, notclosed := <-rmq.connCloseError:
			fmt.Printf("Connection's gone: %v connCloseError len: %v\n",
				err, len(rmq.connCloseError))

			// channel closed by library
			if !notclosed {
				rmq.connCloseError = make(chan *amqp.Error)
				rmq.startConnection <- struct{}{}
			}

		case val, _ := <-rmq.chanCancelError:
			// Will be notified iff queue master node is dead due to
			// channel.Consume is called.
			// i.e publisher side will not be notified with this error.

			// chanCancelError is closed by the library.
			fmt.Printf("Channel's gone :%v chanCancelError len: %v\n",
				val, len(rmq.chanCancelError))

			rmq.chanCancelError = make(chan string)
			rmq.startConnection <- struct{}{}
		case ret, _ := <-rmq.chanReturnError:
			// chanReturnError is closed by the library.
			fmt.Printf("Channel's Publish Return error:%v "+
				"chanReturnError len: %v\n",
				ret, len(rmq.chanReturnError))

			rmq.chanReturnError = make(chan amqp.Return)
			rmq.startConnection <- struct{}{}
		}
	}
}

func consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-rmq.canConsume:
			if rmq.consumeCallbackFunc != nil {
				rmq.consumeCallbackFunc(rmq)
			}
		}
	}
}

func createConnect() {
	rmq.reconn.Store(true)

	runner := func() {
		url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
			rmq.user, rmq.passwd, rmq.host, rmq.port)

		fmt.Printf("Opening Connection(host=%s, username=%s, port=%d, queue=%s)\n",
			rmq.host, rmq.user, rmq.port, rmq.queueName)

		connection, err := amqp.Dial(url)

		if err != nil {
			fmt.Printf("Opening Connection(queue=%s) failed: %v\n",
				rmq.queueName, err)

			fmt.Printf("Wait %d seconds to reconnect...\n", rmq.reconnWait)
			time.Sleep(time.Duration(rmq.reconnWait) * time.Second)

			rmq.reconn.Store(false)
			rmq.startConnection <- struct{}{}

			fmt.Printf("Dial error: %v\n", err)
			return
		}

		rmq.Connection = connection
		rmq.Connection.NotifyClose(rmq.connCloseError)
		rmq.startChannel <- struct{}{}

		fmt.Println("Connection ready...")

		rmq.reconn.Store(false)
	}

	go runner()
}

func createChannel() {
	rmq.rechan.Store(true)

	runner := func() {
		channel, err := rmq.Connection.Channel()
		if err != nil {
			// If there's no connection, stops.
			fmt.Printf("Get Channel error: %v\n", err)
			return
		}

		rmq.Channel = channel

		//Get the queue
		_, err = rmq.Channel.QueueDeclare(
			rmq.queueName,
			rmq.durable, // True: the queue will survive a broker restart
			false,       // delete when unused
			false,       // exclusive
			false,       // noWait
			nil,         // arguments
		)

		if err != nil {
			// If queue Master dead and there's no HA, QueueDeclare fails.
			fmt.Printf("QueueDeclare(queue=%s) failed: %v\n", rmq.queueName, err)
			_ = rmq.Channel.Close()

			fmt.Printf("Wait %d seconds to rechannel...\n", rmq.rechanWait)
			time.Sleep(time.Duration(rmq.rechanWait) * time.Second)

			rmq.rechan.Store(false)
			rmq.startChannel <- struct{}{}
			return
		}

		if rmq.exchange != "" {
			// ExchangeDeclare
			err = rmq.Channel.ExchangeDeclare(
				rmq.exchange,
				rmq.exchangeType,
				true,  // durable
				true,  // auto-delete
				false, // internal
				false, // nowait, wait for confirmation from server
				nil)   // amqp.Table

			// Terminate program if ExchangeDeclare failed.
			if err != nil {
				fmt.Printf("ExchangeDeclare error: %v\n", err)
				p, _ := os.FindProcess(os.Getpid())
				_ = p.Signal(syscall.SIGQUIT)
				return
			}

			var pkey string
			if rmq.mode == debugC {
				pkey = rmq.listeningQueue
			} else {
				pkey = rmq.queueName
			}

			// QueueBind
			err = rmq.Channel.QueueBind(
				rmq.queueName, // queue name
				pkey,
				rmq.exchange, // exchange name
				false,        // nowait, wait for confirmation from server
				nil)

			// Terminate program if QueueBind failed.
			if err != nil {
				fmt.Printf("QueueBind error: %v", err)
				p, _ := os.FindProcess(os.Getpid())
				_ = p.Signal(syscall.SIGQUIT)
				return
			}

			// ExchangeBind for 'debug' mode
			if rmq.mode == debugC {
				err = rmq.Channel.ExchangeBind(
					rmq.exchange,
					rmq.listeningQueue,
					rmq.listeningExchange,
					false, // nowait, wait for confirmation from server
					nil)   // amqp.Table

				if err != nil {
					fmt.Printf("ExchangeBind error: %v\n", err)
					p, _ := os.FindProcess(os.Getpid())
					_ = p.Signal(syscall.SIGQUIT)
					return
				}
			}

		}

		// Register callback only after Exchange And Queue binding done.
		rmq.Channel.NotifyCancel(rmq.chanCancelError)
		rmq.Channel.NotifyReturn(rmq.chanReturnError)

		if rmq.mode != publishC {
			rmq.canConsume <- struct{}{}
		}

		fmt.Println("Channel ready...")

		rmq.rechan.Store(false)
		return
	}

	go runner()
}
