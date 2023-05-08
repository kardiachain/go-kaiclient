package rpc

import "github.com/kardiachain/go-kaiclient/rpc/transport"

type Client struct {
	transport transport.Transport
	addr      string
}

func NewClientWithProxy(addr, proxy string) (*Client, error) {
	c := &Client{
		addr: addr,
	}

	t, err := transport.NewTransport(addr, proxy)
	if err != nil {
		return nil, err
	}
	c.transport = t
	return c, nil
}

func NewClient(addr string) (*Client, error) {
	return NewClientWithProxy(addr, "")
}

func (c *Client) Close() error {
	return c.transport.Close()
}

func (c *Client) Call(method string, out interface{}, params ...interface{}) error {
	return c.transport.Call(method, out, params...)
}
