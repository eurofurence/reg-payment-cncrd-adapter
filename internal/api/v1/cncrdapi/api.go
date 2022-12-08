package cncrdapi

import (
	"net/url"
)

// ErrorDto struct for Error
type ErrorDto struct {
	// The time at which the error occurred.
	Timestamp string `json:"timestamp"`
	// An internal trace id assigned to the error. Used to find logs associated with errors across our services. Display to the user as something to communicate to us with inquiries about the error.
	RequestId string `json:"requestid"`
	// A keyed description of the error. We do not write human readable text here because the user interface will be multi language.  At this time, there are these values: - paylink.parse.error (json body parse error) - paylink.data.invalid (field data failed to validate, see details for more information) - paylink.id.notfound (no such paylink number in the Concardis service) - paylink.id.invalid (syntactically invalid paylink id, must be positive integer) - auth.unauthorized (token missing completely or invalid) - auth.forbidden (permissions missing)
	Message string `json:"message"`
	// Optional list of additional details about the error. If available, will usually contain English language technobabble.
	Details url.Values `json:"details,omitempty"`
}

// HealthReportDto struct for HealthReportDto
type HealthReportDto struct {
	// Health status of this service.
	Status string `json:"status"`
}

// PaymentLinkRequestDto struct for PaymentLinkDto
type PaymentLinkRequestDto struct {
	// The badge number of the attendee. Will be used to build appropriate description, referenceId, etc.
	DebitorId uint64 `json:"debitor_id"`
	// The page title to be shown on the payment page.
	AmountDue int64 `json:"amount_due"`
	// Only used in responses. The total amount paid. TODO - is this Cents or Euros?
	Currency string `json:"currency"`
	// The applicable VAT, in percent.
	VatRate float64 `json:"vat_rate"`
}

// PaymentLinkDto struct for PaymentLinkDto
type PaymentLinkDto struct {
	// The page title to be shown on the payment page.
	Title string `json:"title"`
	// The description to be shown on the payment page.
	Description string `json:"description"`
	// Internal reference number for this payment process.
	ReferenceId string `json:"reference_id"`
	// The purpose of this payment process.
	Purpose string `json:"purpose"`
	// The amount to bill for. TODO - is this Cents or Euros?
	AmountDue int64 `json:"amount_due"`
	// Only used in responses. The total amount paid. TODO - is this Cents or Euros?
	AmountPaid int64 `json:"amount_paid"`
	// The currency to use.
	Currency string `json:"currency"`
	// The applicable VAT, in percent.
	VatRate float64 `json:"vat_rate"`
	// The payment link.
	Link string `json:"link"`
}

// WebhookEventDto struct for WebhookEventDto
type WebhookEventDto struct {
	Transaction WebhookEventTransaction `json:"transaction"`
}

// WebhookEventTransaction struct for WebhookEventTransaction
type WebhookEventTransaction struct {
	Id      int64                          `json:"id"` // id of the transaction (not the payment link)
	Invoice WebhookEventTransactionInvoice `json:"invoice"`
}

// WebhookEventTransactionInvoice struct for WebhookEventTransactionInvoice
type WebhookEventTransactionInvoice struct {
	ReferenceId      string `json:"referenceId"`
	PaymentRequestId int64  `json:"paymentRequestId"` // id of the payment link concerned
}
