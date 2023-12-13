package common

import "time"

// Order представляет основную информацию о заказе
type Order struct {
	OrderUID        string    `json:"order_uid" validate:"required"`
	TrackNumber     string    `json:"track_number" validate:"required"`
	Entry           string    `json:"entry" validate:"required"`
	Delivery        Delivery  `json:"delivery" validate:"required"`
	Payment         Payment   `json:"payment" validate:"required"`
	Items           []Item    `json:"items" validate:"required"`
	Locale          string    `json:"locale" validate:"required"`
	CustomerID      string    `json:"customer_id" validate:"required"`
	DeliveryService string    `json:"delivery_service" validate:"required"`
	ShardKey        string    `json:"shardkey" validate:"required"`
	SmID            int       `json:"sm_id" validate:"required"`
	DateCreated     time.Time `json:"date_created" validate:"required" format:"2006-01-02T06:22:19Z"`
	OofShard        string    `json:"oof_shard" validate:"required"`
}

// Delivery представляет информацию о доставке
type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required"`
}

// Payment представляет информацию о платеже
type Payment struct {
	Transaction  string  `json:"transaction" validate:"required"`
	RequestID    string  `json:"request_id"`
	Currency     string  `json:"currency" validate:"required"`
	Provider     string  `json:"provider" validate:"required"`
	Amount       float64 `json:"amount" validate:"gt=0"`
	PaymentDT    int     `json:"payment_dt" validate:"required"`
	Bank         string  `json:"bank" validate:"required"`
	DeliveryCost float64 `json:"delivery_cost" validate:"gt=0"`
	GoodsTotal   float64 `json:"goods_total" validate:"gt=0"`
	CustomFee    float64 `json:"custom_fee" validate:"gte=0"`
}

// Item представляет информацию о товаре в заказе
type Item struct {
	ChrtID      int    `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price" validate:"required"`
	RID         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        int    `json:"sale" validate:"gte=0"`
	Size        string `json:"size" validate:"gte=0"`
	TotalPrice  int    `json:"total_price" validate:"gte=0"`
	NmID        int    `json:"nm_id" validate:"gte=0"`
	Brand       string `json:"brand" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}
