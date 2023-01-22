package main

import (
	"chargebackapp/workflows"
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	signal := workflows.MerchantSubmissionSignal{
		MerchantResponded: true,
		RespondedAt:       time.Now(),
		Message:           "response from merchant",
		Proof:             "this is the proof",
	}

	err = c.SignalWorkflow(context.Background(),
		"MerchantResponse:rohit+customer@cgocashfree.com",
		"27a5d7a0-f164-41c6-9d2a-e3783445b43d",
		workflows.MerchantSubmissionSignalName, signal)
	if err != nil {
		log.Fatalln("Error sending the Signal", err)
		return
	}

}
