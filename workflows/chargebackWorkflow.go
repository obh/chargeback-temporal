package workflows

import (
	"chargebackapp/models"
	"fmt"

	"go.temporal.io/sdk/log"

	"go.temporal.io/sdk/workflow"
)

const (
	TaskQueue = "chargeback-queue"
)

type ChargebackInput struct {
	Chargeback models.ChargebackRequest
	Customer   models.Customer
	Payment    models.Payment
}

type ChargebackState struct {
	MerchantResponded bool
	Chargeback        models.ChargebackRequest
	Customer          models.Customer
	Payment           models.Payment
	Documents         map[string]interface{}
	MessageHistory    []string
}

type ChargebackResult = ChargebackState

type chargebackWorkflow struct {
	ChargebackState
	runId  string
	logger log.Logger
}

func ChargebackWorkflowId(chargebackId int) string {
	return fmt.Sprintf("Chargeback:%d", chargebackId)
}

func newChargebackWorfklow(ctx workflow.Context, state *ChargebackState) *chargebackWorkflow {
	return &chargebackWorkflow{
		ChargebackState: *state,
		runId:           workflow.GetInfo(ctx).WorkflowExecution.RunID,
		logger:          workflow.GetLogger(ctx),
	}
}

func (w *chargebackWorkflow) pushStatus(ctx workflow.Context, status string) error {
	return workflow.UpsertSearchAttributes(
		ctx,
		map[string]interface{}{
			"chargeBackWorkflowStatus": status,
		},
	)
}

func (w *chargebackWorkflow) waitForMerchantResponse(ctx workflow.Context, email string, cb models.ChargebackRequest) (*MerchantResponseResult, error) {
	var r MerchantResponseResult

	ctx = workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID: fmt.Sprintf("MerchantResponse:%s", email),
	})
	consentWF := workflow.ExecuteChildWorkflow(ctx, MerchantResponse, MerchantResponseWorkflowInput{
		PrimaryEmail:      email,
		ChargebackRequest: ChargebackInput{},
	})
	err := consentWF.Get(ctx, &r)
	return &r, err
}

func (w *chargebackWorkflow) reverseFunds(ctx workflow.Context, payment models.Payment) (models.Payment, error) {
	fmt.Println("demo code")
	return payment, nil
}

func (w *chargebackWorkflow) sendDisputeFailedMail(ctx workflow.Context, payment models.Payment, customer models.Customer) error {
	fmt.Println("demo code")
	return nil
}

func ChargebackProcess(ctx workflow.Context, input *ChargebackInput) (*ChargebackResult, error) {
	w := newChargebackWorfklow(
		ctx,
		&ChargebackState{
			Chargeback:     input.Chargeback,
			Customer:       input.Customer,
			Payment:        input.Payment,
			Documents:      make(map[string]interface{}),
			MessageHistory: make([]string, 10),
		})

	response, err := w.waitForMerchantResponse(ctx, input.Customer.Email, input.Chargeback)
	if err != nil {
		return &w.ChargebackState, err
	}
	w.MerchantResponded = response.MerchantResponded
	if !w.MerchantResponded {
		response, err := w.reverseFunds(ctx, input.Payment)
		if err != nil {
			return &w.ChargebackState, err
		}
		return &w.ChargebackState, w.sendDisputeFailedMail(ctx, response, input.Customer)
	}

	return &w.ChargebackState, nil
}
