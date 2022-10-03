package domain

type Product struct {
	ID                      string         `json:"id" csv:"id"`
	Name                    string         `json:"name" csv:"name"`
	Segment                 ProductSegment `json:"segment" csv:"segment"`
	Price                   float64        `json:"price" csv:"price"`
	DiscountPercentMiniCash float64        `json:"discount_percent_minicash" csv:"discount_percent_minicash"`
}

type ProductSegment string

const (
	MDA ProductSegment = "MDA"
	TV  ProductSegment = "TV"
)
