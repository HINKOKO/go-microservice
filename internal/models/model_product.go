package models

type Product struct {
	ProductID string  `gorm:"primaryKey" json:"productId"`
	Name      string  `json:"name"`
	Price     float32 `gorm:"type:numeric" json:"price"`
	VendorID  string  `json:"vendorId`
}
