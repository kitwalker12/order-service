package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"interfaces"
	"resources"

	restful "github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

type WebServiceHandler struct {
	OrderInteractor OrderInteractor
}

type OrderInteractor interface {
	FindOrder(id string) (order resources.Order, err error)
	FindOrderbyExternalId(id string) (order resources.Order, err error)
	CreateOrder(order *resources.Order) (err error)
}

func (w WebServiceHandler) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/orders").
		Doc("Manage Orders").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{order-id}").To(w.FindOrder).
		Doc("Get Order").
		Operation("FindOrder").
		Param(ws.PathParameter("order-id", "Netsuite Order ID").DataType("string")).
		Writes(resources.Order{}))

	ws.Route(ws.GET("/external/{external-id}").To(w.FindOrderbyExternalId).
		Doc("Get Order by External ID").
		Operation("FindOrder").
		Param(ws.PathParameter("external-id", "External Order ID").DataType("string")).
		Writes(resources.Order{}))

	ws.Route(ws.POST("").To(w.CreateOrder).
		Doc("Create Order").
		Operation("CreateOrder").
		Reads(resources.Order{}))

	container.Add(ws)
}

func (w WebServiceHandler) FindOrder(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("order-id")
	order, err := w.OrderInteractor.FindOrder(id)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: Order could not be found.")
		return
	}
	response.WriteEntity(order)
}

func (w WebServiceHandler) FindOrderbyExternalId(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("external-id")
	order, err := w.OrderInteractor.FindOrderbyExternalId(id)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: Order could not be found.")
		return
	}
	response.WriteEntity(order)
}

func (w WebServiceHandler) CreateOrder(request *restful.Request, response *restful.Response) {
	order := new(resources.Order)
	err := request.ReadEntity(order)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	err = w.OrderInteractor.CreateOrder(order)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteHeaderAndEntity(http.StatusCreated, order)
}

func getReceiver() interfaces.RabbitMQ {
	rabbitmq := interfaces.RabbitMQ{
		Uri:          os.Getenv("RABBITMQ_URL"),
		Exchange:     "services_direct",
		ExchangeType: "direct",
		Reliable:     true,
		Processor: resources.OrderInteractor{
			Publisher: interfaces.RabbitMQ{
				Uri:          os.Getenv("RABBITMQ_URL"),
				Exchange:     "services_direct",
				ExchangeType: "direct",
				Reliable:     true,
			},
			Cache: interfaces.Cache{Uri: os.Getenv("REDIS_URL")},
		},
	}
	return rabbitmq
}

func main() {

	go func() { //start listeners in new go-routines
		c, err := getReceiver().Subscribe("order.create")
		if err != nil {
			log.Fatalf("%s", err)
		}
		//listening on queue forever
		select {}

		if err := c.Shutdown(); err != nil {
			log.Fatalf("error during shutdown: %s", err)
		}
	}()

	go func() { //start listeners in new go-routines
		c, err := getReceiver().Subscribe("order.create_multiple_orders")
		if err != nil {
			log.Fatalf("%s", err)
		}
		//listening on queue forever
		select {}

		if err := c.Shutdown(); err != nil {
			log.Fatalf("error during shutdown: %s", err)
		}
	}()

	go func() { //start listeners in new go-routines
		c, err := getReceiver().Subscribe("order.update")
		if err != nil {
			log.Fatalf("%s", err)
		}
		//listening on queue forever
		select {}

		if err := c.Shutdown(); err != nil {
			log.Fatalf("error during shutdown: %s", err)
		}
	}()

	go func() { //start listeners in new go-routines
		c, err := getReceiver().Subscribe("order.pending_user_created")
		if err != nil {
			log.Fatalf("%s", err)
		}
		//listening on queue forever
		select {}

		if err := c.Shutdown(); err != nil {
			log.Fatalf("error during shutdown: %s", err)
		}
	}()

	go func() { //start listeners in new go-routines
		c, err := getReceiver().Subscribe("order.fulfill")
		if err != nil {
			log.Fatalf("%s", err)
		}
		//listening on queue forever
		select {}

		if err := c.Shutdown(); err != nil {
			log.Fatalf("error during shutdown: %s", err)
		}
	}()

	go func() { //start listeners in new go-routines
		c, err := getReceiver().Subscribe("return.pending.netsuite_order_id_required")
		if err != nil {
			log.Fatalf("%s", err)
		}
		//listening on queue forever
		select {}

		if err := c.Shutdown(); err != nil {
			log.Fatalf("error during shutdown: %s", err)
		}
	}()

	go func() { //start listeners in new go-routines
		c, err := getReceiver().Subscribe("order.service_status")
		if err != nil {
			log.Fatalf("%s", err)
		}
		//listening on queue forever
		select {}

		if err := c.Shutdown(); err != nil {
			log.Fatalf("error during shutdown: %s", err)
		}
	}()

	wsContainer := restful.NewContainer()
	w := WebServiceHandler{
		OrderInteractor: resources.OrderInteractor{
			Publisher: interfaces.RabbitMQ{
				Uri:          os.Getenv("RABBITMQ_URL"),
				Exchange:     "services_direct",
				ExchangeType: "direct",
				Reliable:     true,
			},
			Cache: interfaces.Cache{Uri: os.Getenv("REDIS_URL")},
		},
	}
	w.Register(wsContainer)

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:3000/apidocs and enter http://localhost:3000/apidocs.json in the api input field.
	url := fmt.Sprintf("%s://%s%s", os.Getenv("PROTO"), os.Getenv("DOCKERCLOUD_SERVICE_FQDN"), os.Getenv("PORT"))
	config := swagger.Config{
		WebServices:    wsContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: url,
		ApiPath:        "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/go/src/swagger-ui/dist"}
	swagger.RegisterSwaggerService(config, wsContainer)

	log.Printf("start listening")
	server := &http.Server{Addr: os.Getenv("PORT"), Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
