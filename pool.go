package pool

import (
	"errors"
	"reflect"
	"strings"
	"time"
)

type pool struct {
	objectType    string
	capacity      int
	minimumObject int
	idleObject    chan abstractObjectInterface
	runningObject int
	objectCreator objectFactory
	timeOut       time.Duration
}

//NewEmptyPool creates a new empty pool
func NewEmptyPool() *pool {
	return &pool{}
}

func (p *pool) isObjectTypeConfiged() bool {
	if p.objectType == "" || p.objectCreator == nil {
		return false
	}
	return true
}

//Config the type of object created from the pool
func (p *pool) SetObjectType(objectType string) error {
	if p.isObjectTypeConfiged() {
		return errors.New("Pool has been configed")
	}

	factory, supported := supportedObject[strings.ToLower(objectType)]
	if !supported {
		return errors.New("Unsupported object type")
	}

	p.objectType = strings.ToLower(objectType)
	p.objectCreator = factory
	return nil
}

//Config the maximum and minimum numbers of objects in the pool
func (p *pool) SetupObjectPool(capacity int, minObject int) error {
	if capacity < minObject {
		return errors.New("Mininum object number must be larger than pool's capacity")
	}
	if !p.isObjectTypeConfiged() {
		return errors.New("Please config object type")
	}

	p.timeOut = 5 * time.Second //config the default timeout = 5 seconds
	p.capacity = capacity
	p.minimumObject = minObject
	p.idleObject = make(chan abstractObjectInterface, capacity)
	for i := 0; i < minObject; i++ {
		newObject, err := p.objectCreator.createObject()
		if err != nil {
			return err
		}
		p.idleObject <- newObject
	}
	return nil
}

//Config the live-time for every object created from the pool
func (p *pool) SetObjectTimeOut(long time.Duration) {
	p.timeOut = long
}

//Get the number of running objects in the pool
func (p *pool) GetRunningObjectNumber() int {
	return p.runningObject
}

//Get the number of idle objects in the pool
func (p *pool) GetIdleObjectNumber() int {
	return len(p.idleObject)
}

//Release an object from its pool
func (p *pool) GetObjectFromPool() (abstractObjectInterface, error) {
	if !p.isObjectTypeConfiged() {
		return nil, errors.New("Please config object type")
	}

	if len(p.idleObject) != 0 {
		object := <-p.idleObject
		p.runningObject++
		return object, nil
	}

	if p.runningObject < p.capacity {
		object, err := p.objectCreator.createObject()
		if err != nil {
			return nil, err
		}
		p.runningObject++
		return object, nil
	}
	return nil, errors.New("Pool is full")
}

//Return an object to its pool
func (p *pool) ReturnObjectToPool(object abstractObjectInterface) error {
	objectType := strings.ToLower(reflect.TypeOf(object).Name())
	_, supported := supportedObject[objectType]
	if !supported {
		return errors.New("Unknown object")
	}

	p.idleObject <- object
	p.runningObject--
	return nil

}

//Remove timeout objects and create the new ones to achive the minimum number
func (p *pool) RefreshPool() error {
	//remove out-of-time objects
	for i := 0; i < len(p.idleObject); i++ {
		object := <-p.idleObject
		if object.getAliveTime() < p.timeOut {
			p.idleObject <- object
		}
	}

	//create new objects to achive the minimum number
	totalObject := len(p.idleObject) + p.runningObject
	for ; totalObject < p.minimumObject; totalObject++ {
		newObject, err := p.objectCreator.createObject()
		if err != nil {
			return err
		}
		p.idleObject <- newObject
	}
	return nil
}

var supportedObject = map[string]objectFactory{
	"connection": &connectionFactory{},
	"pencil":     &pencilFactory{},
}
