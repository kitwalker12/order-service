package resources

type OrderInteractor struct {
	Order Order
}

type Order struct {
}

func (o OrderInteractor) FindOrder(id string) (order Order, err error) {
	return Order{}, nil
}

func (o OrderInteractor) CreateOrder(order *Order) (err error) {
	return nil
}
