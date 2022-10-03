package domain

type LegalDocument struct {
	DocumentType   string `json:"document_type"`
	DocumentNumber string `json:"document_number"`
}
type UserProfile struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	LegalDocument LegalDocument `json:"legal_document"`
	AnnualIncome  int           `json:"annual_income"`
	CreditLimit   float64       `json:"credit_limit"`
}
type EWallet struct {
	CreditLimit   float64 `json:"credit_limit"`
	UsedCredit    float64 `json:"used_credita"`
	LoyaltyPoints int     `json:"loyalty_points"`
}
