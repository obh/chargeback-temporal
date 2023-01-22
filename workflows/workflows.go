package workflows

import "time"

type MerchantResponseWorkflowInput struct {
	PrimaryEmail      string
	ChargebackRequest ChargebackInput
}

type MerchantSubmission struct {
	MerchantResponded bool
	RespondedAt       time.Time
	Message           string
	Proof             string
}

type MerchantSubmissionSignal = MerchantSubmission

// type MerchantProof struct {
// 	Message string
// 	Files   []string
// }
