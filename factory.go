package pool

import (
	"net"
	"time"
)

type objectFactory interface {
	createObject() (abstractObjectInterface, error)
}

type connectionFactory struct{}

func (c *connectionFactory) createObject() (abstractObjectInterface, error) {
	var connection connection
	con, err := net.Dial("tcp", "localhost:1433")
	if err != nil {
		return nil, err
	}

	connection.connector = con
	connection.createAt = time.Now()
	return connection, nil
}

type pencilFactory struct{}

func (p *pencilFactory) createObject() (abstractObjectInterface, error) {
	return pencil{}, nil
}
