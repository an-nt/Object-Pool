package pool

import "net"

type Connection struct {
	Connector net.Conn
}

type Pencil struct {
	Brand string
}
