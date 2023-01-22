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

	// err := w.pushStatus(ctx, "pending_merchant_response")
	// if err != nil {
	// 	return &r, err
	// }

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

	return &w.ChargebackState, nil
}
