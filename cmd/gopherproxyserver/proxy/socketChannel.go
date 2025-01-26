package proxy

type SocketChannel struct {
	Id          string
	Source      *Client
	Sink        *Client
	Initialized bool
}
