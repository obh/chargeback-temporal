package workflows

import "time"

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
