package pool

import "net"

type ObjectFactory interface {
	CreateObject() (interface{}, error)
}

type connectionFactory struct{}

func (c *connectionFactory) CreateObject() (interface{}, error) {
	con, err := net.Dial("tcp", "localhost:1433")
	if err != nil {
		return nil, err
	}

	result := &Connection{
		Connector: con,
	}
	return result, nil

}

type pencilFactory struct{}

func (p *pencilFactory) CreateObject() (interface{}, error) {
	return &Pencil{}, nil
}
