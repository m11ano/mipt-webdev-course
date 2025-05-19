package productsw

import (
	tclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const (
	ProductsQueue = "products-queue"
)

type ProductsWorker worker.Worker

func NewWorker(tClient tclient.Client) ProductsWorker {
	w := worker.New(tClient, ProductsQueue, worker.Options{})

	return w
}
