package go253

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/ratelimit"
)

type Client interface {
	// In() send-only channel for messages
	In() chan<- *message
	// Out() read-only channel for send records
	Out() <-chan SendRecord
	// Close() close the client, close the channels
	Close()
	// Balance() get the balance of the account
	Balance() (float64, error)
	// SingleMessage() send a single message
	SingleMessage(msg *message) SendRecord
}

type client struct {
	ctx        context.Context    // context
	cancel     context.CancelFunc // cancel function
	limiter    ratelimit.Limiter  // ratelimit
	account    string             // 253 account
	password   string             // 253 password
	nodeType   int                // node number
	msgChan    chan *message      // channel for messages
	recordChan chan SendRecord    // channel for send records
}

// validate validates the client.
func (c *client) validate() error {
	if c.account == "" || c.password == "" {
		return errors.New("account and password cannot be empty")
	}
	if _, ok := nodeName[c.nodeType]; !ok {
		return errors.New("invalid node setting")
	}
	return nil
}

func NewClient(ops ...clientOption) (Client, error) {
	c := &client{}
	for _, op := range ops {
		op(c)
	}
	if err := c.validate(); err != nil {
		return nil, errors.Wrap(err, "client validation failed")
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	if c.limiter == nil {
		c.limiter = ratelimit.NewUnlimited()
	}
	return c, nil
}

func (c *client) In() chan<- *message {
	c.msgChan = make(chan *message, 1)
	c.recordChan = make(chan SendRecord, 1)
	c.send()
	return c.msgChan
}

func (c *client) Out() <-chan SendRecord {
	return c.recordChan
}

func (c *client) Close() {
	defer close(c.msgChan)
	c.cancel()
}

func (c *client) send() {
	go func() {
		defer close(c.recordChan)
		for {
			select {
			case <-c.ctx.Done():
				return
			case m, ok := <-c.msgChan:
				if !ok {
					return
				}
				m.Account = c.account
				m.Password = c.password
				c.limiter.Take() // ratelimit
				go c.postMessage(m)
			default:
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()
}
