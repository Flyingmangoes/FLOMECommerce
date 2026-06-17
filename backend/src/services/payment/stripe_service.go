package payment_service

import (
	"backend/src/config"
	Logger "backend/src/utils/logger"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/stripe/stripe-go/v85"
	"go.uber.org/zap"
)

type PaymentService struct {
	WebhookKey string
	SuccessURL string
	CancelURL string
	StripeClient *stripe.Client
}

type ItemDetail struct {
    ProductID   string
    ProductName string
    ProductDesc string
    StoreName   string
    Price       float64
    Quantity    int
    ImageUrl    string
}

func NewPaymentService(cfg *config.StripeConfig, successURL, cancelURL string) *PaymentService {
	return &PaymentService{
		SuccessURL: successURL,
		CancelURL: cancelURL,
		WebhookKey: cfg.STRIPE_WEBHOOK_SECRET,
		StripeClient: stripe.NewClient(cfg.STRIPE_SECRET_KEY),
	}
}

func NewStripeURL(cfg *config.ConfigManager)([]string){
	u := make([]string,0)
	fhp := net.JoinHostPort(cfg.SERV_CONF.FrontendHOST, cfg.SERV_CONF.FrontendPORT)

	success := fmt.Sprintf("http://%s/success", fhp)
	cancel := fmt.Sprintf("http://%s/cancel", fhp)

	u = append(u, success, cancel)
	return u
}

func(ps *PaymentService) CreateCheckoutSession(ctx context.Context, orderID, buyerEmail string, items []ItemDetail) (*stripe.CheckoutSession, error) {
	lineItems := make([]*stripe.CheckoutSessionCreateLineItemParams, 0, len(items))

	for _, item := range items{	
		priceInCents := int64(item.Price * 100)

		lineItems = append(lineItems, &stripe.CheckoutSessionCreateLineItemParams{
			Quantity: stripe.Int64(int64(item.Quantity)),
			PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
				Currency: stripe.String("usd"),
				ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
					Name: &item.ProductName,
					Description: &item.ProductDesc,
					Images: []*string{stripe.String(item.ImageUrl)},
					Metadata: map[string]string{
						"product_id": item.ProductID,
						"store_name": item.StoreName,
					},
				},
				UnitAmount: stripe.Int64(priceInCents),
			},
		})
	}

	success := fmt.Sprintf("http://%s?session_id={CHECKOUT_SESSION_ID}", ps.SuccessURL)
	cancel := fmt.Sprintf("http://%s", ps.CancelURL)

	Logger.Log.Debug("urls", zap.String("success url", success))
	Logger.Log.Debug("urls", zap.String("cancel url", cancel))

	params := &stripe.CheckoutSessionCreateParams{
		LineItems: lineItems,
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		CustomerEmail: stripe.String(buyerEmail),
		SuccessURL: stripe.String(success),
		CancelURL: stripe.String(cancel),
		ExpiresAt: stripe.Int64(time.Now().Add(30 * time.Minute).Unix()),

		Metadata: map[string]string{
			"order_id": orderID,
		},
	}

	s, err := ps.StripeClient.V1CheckoutSessions.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	return s, nil
	
}

