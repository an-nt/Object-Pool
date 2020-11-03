package pool

import (
	"errors"
	"reflect"
	"strings"
)

type pool struct {
	objectType    string
	idleObject    chan interface{}
	runningObject int
	capacity      int
	objectCreator ObjectFactory
}

func NewEmptyPool() *pool {
	return &pool{}
}
func (p *pool) isObjectTypeConfiged() bool {
	if p.objectType == "" || p.objectCreator == nil {
		return false
	}
	return true
}

func (p *pool) SetObjectType(objectType string) error {
	if p.isObjectTypeConfiged() {
		return errors.New("Pool has been configed")
	}

	factory, supported := supportedObject[strings.ToLower(objectType)]
	if !supported {
		return errors.New("Unsupported object type")
	}
	//fmt.Println(reflect.TypeOf(factory).Name())
	p.objectType = strings.ToLower(objectType)
	p.objectCreator = factory
	return nil
}

func (p *pool) SetupObjectPool(capacity int, minObject int) error {
	if capacity < minObject {
		return errors.New("Mininum object number must be larger than pool's capacity")
	}
	if !p.isObjectTypeConfiged() {
		return errors.New("Please config object type")
	}
	p.capacity = capacity
	p.idleObject = make(chan interface{}, capacity)
	for i := 0; i < minObject; i++ {
		newObject, err := p.objectCreator.CreateObject()
		if err != nil {
			return err
		}
		p.idleObject <- newObject
	}
	return nil
}

func (p *pool) GetObjectFromPool() (interface{}, error) {
	if !p.isObjectTypeConfiged() {
		return nil, errors.New("Please config object type")
	}

	if len(p.idleObject) != 0 {
		object := <-p.idleObject
		p.runningObject++
		return object, nil
	}

	if p.runningObject < p.capacity {
		object, err := p.objectCreator.CreateObject()
		if err != nil {
			return nil, err
		}
		p.runningObject++
		return object, nil
	}
	return nil, errors.New("Pool is full")
}

func (p *pool) ReturnObjectToPool(object interface{}) error {
	objectType := strings.ToLower(reflect.TypeOf(object).Name())
	_, supported := supportedObject[objectType]
	if !supported {
		return errors.New("Unknown object")
	}

	p.idleObject <- object
	p.runningObject--
	return nil

}

var supportedObject = map[string]ObjectFactory{
	"connection": &connectionFactory{},
	"pencil":     &pencilFactory{},
}
