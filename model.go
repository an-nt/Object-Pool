package pool

import (
	"net"
	"time"
)

type abstractObjectInterface interface {
	getAliveTime() time.Duration
}

type abstractObject struct {
	createAt time.Time
}

func (a abstractObject) getAliveTime() time.Duration {
	return time.Since(a.createAt)
}

type connection struct {
	abstractObject
	connector net.Conn
}

type pencil struct {
	abstractObject
	Brand string
}
