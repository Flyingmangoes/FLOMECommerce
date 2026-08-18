package models

type Cart struct {
	ID 			string `db:"cart_id"`
	UserID 		string `db:"user_id"`
	CartItems 	[]CartItem	
}

type CartItem struct {
	ID			string `db:"cart_item_id"`
	CartID 		string `db:"cart_id"`
	ProductID 	string `db:"product_id"`
	Quantity 	int	   `db:"quantity"`
}