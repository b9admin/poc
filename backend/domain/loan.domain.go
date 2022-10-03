package domain

import (
	"cloud.google.com/go/civil"
)

type LoanApplication struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	SalesOrder       SalesOrder `json:"sales_order"`
	TotalPayment     float64    `json:"total_payment"`
	StartDate        civil.Date `json:"start_date"`
	TenorMonths      int        `json:"tenor_months"`
	EndDate          civil.Date `json:"end_date"`
	InstallmentValue float64    `json:"installment_value"`
	RateOfInterest   float64    `json:"rate_of_interest"`
}

type LoanPayment struct {
	ID                string  `json:"id"`
	LoanApplicationID string  `json:"loan_application_id"`
	Payment           float64 `json:"payment"`
	Date              float64 `json:"date"`
	PaymentMade       bool    `json:"payment_made"`
}

const (
	PAYMENT_SUCCESS string = "PAYMENT_SUCCESS"
	PAYMENT_FAILED  string = "PAYMENT_FAILED"
	PAYMENT_PARTIAL string = "PAYMENT_PARTIAL"
)

type InterestApportionment struct {
	CustomerSegment string  `json:"customer_segment" csv:"customer_segment"`
	Tenor           int     `json:"tenor" csv:"tenor"`
	ProductSegment  string  `json:"product_segment" csv:"product_segment"`
	Percent         float64 `json:"percent" csv:"percent"`
}
