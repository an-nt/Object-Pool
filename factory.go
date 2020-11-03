package pool

import (
	"fmt"
	"net"
)

type ObjectFactory interface {
	CreateObject() (interface{}, error)
}

type connectionFactory struct{}

func (c *connectionFactory) CreateObject() (interface{}, error) {
	con, err := net.Dial("tcp", "localhost:1433")
	if err != nil {
		return nil, err
	}

	fmt.Println("Create a new connection")
	return Connection{Connector: con}, nil
}

type pencilFactory struct{}

func (p *pencilFactory) CreateObject() (interface{}, error) {
	return Pencil{}, nil
}
