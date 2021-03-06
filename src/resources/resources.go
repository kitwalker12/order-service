package resources

import "log"

type OrderInteractor struct {
	Order     Order
	Publisher Publisher
	Cache     Cache
}

type Order struct {
}

type Publisher interface {
	Publish(routingKey string, body string) error
}

type Cache interface {
	Do(cmd string, args ...interface{}) (interface{}, error)
}

func (o OrderInteractor) FindOrder(id string) (order Order, err error) {
	return Order{}, nil
}

func (o OrderInteractor) FindOrderbyExternalId(id string) (order Order, err error) {
	return Order{}, nil
}

func (o OrderInteractor) CreateOrder(order *Order) (err error) {
	return nil
}

func (o OrderInteractor) Process(routingKey string, body []byte) {
	log.Printf("Processing %s for %s", string(body), routingKey)
}
