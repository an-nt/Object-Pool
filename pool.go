package pool

import (
	"errors"
	"reflect"
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

func (p *pool) SetObjectType(objectType string) (*pool, error) {
	if p.isObjectTypeConfiged() {
		return p, errors.New("Pool has been configed")
	}

	factory, supported := supportedObject[objectType]
	if !supported {
		return p, errors.New("Unsupported object type")
	}

	p.objectType = objectType
	p.objectCreator = factory
	return p, nil
}

func (p *pool) SetupObjectPool(capacity int, minObject int) (*pool, error) {
	if capacity < minObject {
		return p, errors.New("Mininum object number must be larger than pool's capacity")
	}
	if !p.isObjectTypeConfiged() {
		return p, errors.New("Please config object type")
	}

	p.idleObject = make(chan interface{}, capacity)
	for i := 0; i < minObject; i++ {
		newObject, err := p.objectCreator.CreateObject()
		if err != nil {
			return p, err
		}
		p.idleObject <- newObject
		p.runningObject++
	}
	return p, nil
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
	reflecter := reflect.TypeOf(object)
	objectName := reflecter.Name()
	_, supported := supportedObject[objectName]
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
