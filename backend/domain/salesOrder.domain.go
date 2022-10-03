package domain

type LineItem struct {
	Product  Product `json:"product"`
	Quantity int32   `json:"quantity"`
}
type SalesOrder struct {
	ID              string     `json:"id"`
	LineItems       []LineItem `json:"line_items"`
	ShippingAddress string     `json:"shipping_address"`
	TotalOrderPrice float64    `json:"total_order_price"`
}
