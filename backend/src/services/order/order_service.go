package order_services

import (
	"backend/src/models"
	"backend/src/repository"
	"backend/src/services"
	"context"
	"database/sql"
	"fmt"
)

type OrderService struct{
	Tx *services.TxManager
}

type PlaceOrderParams struct {
    BuyerID     		string
    BuyerEmail  		string
    CombinedLocation 	string
    Status      		string
    ProductList 		[]repository.OrderItemInput
}

type CancelOrderParams struct {
	OrderId 			string
	BuyerId 			string
}

func (os *OrderService) PlaceOrder(ctx context.Context, params *PlaceOrderParams) (*models.Order, error) {
	var results *models.Order

	err := os.Tx.WithTx(ctx, func(tx *sql.Tx) error {
		products, err := os.Tx.Products.GetProductForUpdate(ctx, tx, params.ProductList)
		if err != nil { return err }

		for _, p := range products {
			for _, req := range params.ProductList {
				if req.ProductID == &p.ProductID && p.Availability < *req.Quantity {
					return fmt.Errorf("insufficient stock: %s", p.ProductID)
				}
			}
		}

		total := calculateTotal(products, params.ProductList)

		orderItems := make([]repository.OrderItemInput, 0)
		for _, p := range params.ProductList {
			for _, prod := range products {
				if prod.ProductID == *p.ProductID {
					orderItems = append(orderItems, repository.OrderItemInput{
						ProductID: &prod.ProductID,
						Price: &prod.Price,
						Quantity: p.Quantity,
					})
				}
			}
		}

		order, _, err := os.Tx.Orders.CreateOrder(ctx, tx, &repository.OrderStoreParams{
			BaseParams: repository.BaseParams{
				UserId: &params.BuyerID,
				Email:  &params.BuyerEmail,
			},

            TotalPrice:  &total,
            Location:    &params.CombinedLocation,
            Status:      &params.Status,
            ProductList: orderItems,
		})

		if err != nil { return err }

		if err := os.Tx.Products.DeductStock(ctx, tx, orderItems); err != nil {
			return err
		}

		results = order
		return nil
	})

	return results, err
}

func (os *OrderService)CancelOrder(ctx context.Context, params *CancelOrderParams) error {
	err := os.Tx.WithTx(ctx, func(tx *sql.Tx) error {
		err := os.Tx.Orders.RemoveOrder(ctx, tx, &repository.OrderStoreParams{
			BaseParams: repository.BaseParams{ UserId: &params.BuyerId },
			OrderID: &params.OrderId,
		})

		if err != nil { return err }

		return nil
	})

	if err != nil {
		return err 
	}

	return nil
}

func calculateTotal(products []models.Product, productList []repository.OrderItemInput) float64 {
	priceMap := make(map[string]float64)
	for _, p := range products {
		priceMap[p.ProductID] = p.Price
	}

	var total float64
	for _, item := range productList {
		total += priceMap[*item.ProductID] * float64(*item.Quantity)
	}

	return total
}