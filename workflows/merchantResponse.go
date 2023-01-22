package workflows

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	MerchantSubmissionSignalName = "merchant-submission"
	MerchantSubmissionPeriod     = time.Hour * 24 * 7
)

type MerchantResponseResult struct {
	MerchantResponded bool
	RespondedAt       time.Time
	Message           string
	Proof             string
}

type SendEmailInput struct {
	Email   string
	message string
}
type SendEmailResult struct {
	Status bool
}

func SendEmail(input SendEmailInput) (*SendEmailResult, error) {
	var result SendEmailResult

	fmt.Println("Sending email to merchant")
	return &result, nil
}

func emailMerchant(ctx workflow.Context, input *MerchantResponseWorkflowInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	i := SendEmailInput{Email: input.PrimaryEmail, message: "Some standard text"}
	f := workflow.ExecuteActivity(ctx, SendEmail, i)

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

func MerchantResponse(ctx workflow.Context, input *MerchantResponseWorkflowInput) (*MerchantResponseResult, error) {
	err := emailMerchant(ctx, input)
	if err != nil {
		return &MerchantResponseResult{}, err
	}
	submission, err := waitForSubmission(ctx)

	result := MerchantResponseResult(*submission)
	return &result, err
}
