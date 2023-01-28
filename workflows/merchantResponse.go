package workflows

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	MerchantSubmissionSignalName = "merchant-submission"
	MerchantSubmissionPeriod     = time.Minute * 15
)

type MerchantResponseResult struct {
	MerchantResponded bool
	RespondedAt       time.Time
	Message           string
	Proof             string
}

type SendEmailResult struct {
	Status bool
}

func InvokeNotifyAPI(input ChargebackWFInput) (SendEmailResult, error) {
	var result SendEmailResult
	body := map[string]uint{
		"payment_id":    input.Payment.ID,
		"chargeback_id": input.Chargeback.ID,
	}
	fmt.Println("sending request to InvokeNotifyAPI:", body)
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(body)
	_, err := http.Post("http://localhost:1323/notify", "application/json", buffer)
	if err != nil {
		return result, err
	}

	result.Status = true
	return result, nil
}

func emailMerchant(ctx workflow.Context, input *ChargebackWFInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	// i := SendEmailInput{Email: input.Merchant.PrimaryEmail, message: "Some standard text"}
	f := workflow.ExecuteActivity(ctx, InvokeNotifyAPI, *input)

	return f.Get(ctx, nil)
}

func waitForSubmission(ctx workflow.Context) (*MerchantSubmission, error) {
	var response MerchantSubmission
	var err error

	s := workflow.NewSelector(ctx)

	ch := workflow.GetSignalChannel(ctx, MerchantSubmissionSignalName)
	s.AddReceive(ch, func(c workflow.ReceiveChannel, more bool) {
		var submission MerchantSubmissionSignal
		c.Receive(ctx, &submission)

		response = MerchantSubmission(submission)
	})
	s.AddFuture(workflow.NewTimer(ctx, MerchantSubmissionPeriod), func(f workflow.Future) {
		err = f.Get(ctx, nil)
		response.MerchantResponded = true
	})
	s.Select(ctx)

	return &response, err
}

// func MerchantResponse(input ChargebackWFInput) (*MerchantResponseResult, error) {
// 	err := emailMerchant(input)
// 	if err != nil {
// 		return &MerchantResponseResult{}, err
// 	}
// 	submission, err := waitForSubmission(ctx)

// 	result := MerchantResponseResult(*submission)
// 	return &result, err
// }
