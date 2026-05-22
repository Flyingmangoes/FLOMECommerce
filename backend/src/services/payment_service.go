package services

import (
	"backend/src/config"
	"backend/src/repository"
	"backend/src/utils"
	Logger "backend/src/utils/logger"
	"context"
	"net"
	"time"

	"github.com/stripe/stripe-go/v85"
	"go.uber.org/zap"
)

type PaymentService struct {
 	StripeClient *stripe.Client
	Product_repo repository.ProductStoreInterface
	Domain string
}

type ItemDetail struct {
	ProductID 	string 
	ProductName string 
	ProductDesc string
	StoreName 	string 
    Quantity  	int    
}

type OrderDetail struct {
	OrderID 	string
    BuyerID     string      
	BuyerEmail 	string
    Location    string     
}

func SetupPayment(cfg *config.Application) *PaymentService {
	return &PaymentService{
		StripeClient: stripe.NewClient(cfg.STRIPE_CONF.STRIPE_PUBLIC_KEY),
		Domain: net.JoinHostPort(cfg.SERV_CONF.FrontendHOST, cfg.SERV_CONF.FrontendPORT),
	}
}

func(payment *PaymentService) CreateCheckoutSession(ctx context.Context, items *[]ItemDetail, details *[]OrderDetail) (*stripe.CheckoutSession, error) {
	lineItems := make([]*stripe.CheckoutSessionCreateLineItemParams, 0, len(*items))
	var params *stripe.CheckoutSessionCreateParams

	for _, item := range *items{
		price, err := payment.Product_repo.GetProductByID(ctx, item.ProductID)
		
		if err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			return nil, err
		}

		lineItems = append(lineItems, &stripe.CheckoutSessionCreateLineItemParams{
			Quantity: stripe.Int64(int64(item.Quantity)),
			PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
				Currency: stripe.String("usd"),
				ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
					Name: &item.ProductName,
					Description: &item.ProductDesc,
					Metadata: map[string]string{
						"product_id": item.ProductID,
						"store_name": item.StoreName,
					},
				},
				UnitAmount: stripe.Int64(int64(price.Price)),
			},
		})
	}

	for i := 0; i < len(*details); i++ {
		params = &stripe.CheckoutSessionCreateParams{
			LineItems: lineItems,
			Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
			SuccessURL: stripe.String(payment.Domain + "?success=true"),
			CancelURL: stripe.String(payment.Domain + "?canceled=true"),

			Metadata: map[string]string{
				"order_id": (*details)[i].OrderID,
				"buyer_id": (*details)[i].BuyerID,
			},

			CustomerEmail: utils.STRtoptr((*details)[i].BuyerEmail),
			ExpiresAt: utils.I64toptr(time.Now().Add(30 * time.Minute).Unix()),
		}
	}

	sc, err := payment.StripeClient.V1CheckoutSessions.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	return sc, nil
}

