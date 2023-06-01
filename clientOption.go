package go253

import (
	"time"

	"go.uber.org/ratelimit"
)

type clientOption func(*client)

// WithAccount sets the account for the client.
func WithAccount(a string) clientOption {
	return func(c *client) {
		c.account = a
	}
}

// WithPassword sets the password for the client.
func WithPassword(p string) clientOption {
	return func(c *client) {
		c.password = p
	}
}

// WithNodeType sets the node type for the client.
func WithNodeType(nodeType int) clientOption {
	return func(c *client) {
		if _, ok := nodeName[nodeType]; ok {
			c.nodeType = nodeType
		}
	}
}

// WithLimit sets the ratelimit for the client.
// param limit: the ratelimit number; param duration: the ratelimit duration.
func WithLimit(limit int, duration time.Duration) clientOption {
	return func(c *client) {
		if limit > 0 && duration > 0 {
			c.limiter = ratelimit.New(limit, ratelimit.WithoutSlack, ratelimit.Per(duration))
		}
	}
}
