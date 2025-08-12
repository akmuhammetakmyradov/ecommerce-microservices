package models

type CartItem struct {
	UserID int64
	SKU    uint32
	Count  uint32
}

type DeleteCartItem struct {
	UserID int64
	SKU    uint32
}

type CartItemModel struct {
	SKU   uint32
	Count uint32
	Name  string
	Price uint32
}

type CartItemsList struct {
	Items      []CartItemModel
	TotalPrice uint32
}

type StockItem struct {
	SKU      uint32
	Name     string
	Type     string
	Count    uint32
	Price    uint32
	Location string
}
