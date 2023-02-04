package main

import (
	"log"

	"github.com/obh/chargebackapp/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: "temporaltest-frontend:7233",
	})
	if err != nil {
		log.Println("this is a new client")
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, workflows.TaskQueue, worker.Options{})

	w.RegisterWorkflow(workflows.ChargebackProcess)
	// w.RegisterActivity(workflows.MerchantResponse)
	w.RegisterActivity(workflows.InvokeNotifyAPI)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
