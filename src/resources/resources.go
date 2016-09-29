package resources

import "log"

type OrderInteractor struct {
	Order     Order
	Publisher Publisher
}

type Order struct {
}

type Publisher interface {
	Publish(routingKey string, body string) error
}

func (o OrderInteractor) FindOrder(id string) (order Order, err error) {
	err = o.Publisher.Publish("order.find", id)
	return Order{}, err
}

func (o OrderInteractor) CreateOrder(order *Order) (err error) {
	return nil
}

func (o OrderInteractor) Process(routingKey string, body []byte) {
	log.Printf("Processing %s for %s", string(body), routingKey)
}
