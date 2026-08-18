package email_services

import (
	"backend/src/config"
	"backend/src/repository"
	logger_system "backend/src/utils/LoggerSystem"
	"fmt"
	"strconv"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

type MailType int

type MailService interface {
	CreateMail(mailReq *Mail) []byte
	SendMail(mailReq *Mail) error
}

type MailServiceParams struct {
	From 		string
	To 			[]string
	Subject 	string
	MailType 	MailType
	MailData 	*MailData
}

type SendgridManager struct {
	Users 			repository.UserStoreInterface
	SendgridEmail
	SendgridContext
	SendgridTemplate
}	

type SendgridEmail struct {
	CustomerSupportEmail 	string
	DomainEmail 			string
	TestEmail 				string
}

type SendgridTemplate struct{
	TUserVerification 		string
	TPassReset 				string
	TOrderConfirmation 		string
}

type SendgridContext struct {
	ServerUrl 					string
	SendgridApiSecret 			string
	UserVerificationSecret 		string
}

func NewSendgridManager(cfg *config.SendgridConfig, server_url string, us repository.UserStoreInterface) *SendgridManager {
	return &SendgridManager{
		Users: us,
		SendgridEmail: SendgridEmail{
			CustomerSupportEmail: cfg.SUPPORT_EMAIL,
			DomainEmail:cfg.DOMAIN_EMAIL,
			TestEmail: cfg.TEST_EMAIL,
		},
		SendgridContext: SendgridContext{
			ServerUrl: server_url,
			SendgridApiSecret: cfg.SENDGRID_SECRET,
			UserVerificationSecret: cfg.USER_VERIFICATION_SECRET,
		},
		SendgridTemplate: SendgridTemplate{
			TUserVerification: cfg.TEMPLATE_USER_VERIFICATION,
			TPassReset: cfg.TEMPLATE_PASS_RESET,
			TOrderConfirmation: cfg.TEMPLATE_ORDER_CONFIRMATION,
		},
	}
}

type MailData struct {
	DomainName 		string 		`json:"domainName"`
	Username 		string		`json:"username"`
	FirstName		string		`json:"firstName"`
	SupportEmail 	string		`json:"supportEmail"`
	TokenUrl 		string		`json:"token"`
	Expiration 		*time.Time 	`json:"expires"`
}

type Mail struct {
	from 		string
	to 			[]string
	subject 	string
	mailTyp 	MailType
	data 		*MailData
}

const (
	UserVerification MailType = iota + 1
	PassReset 
	OrderConfirmation
)

var typeList = map[MailType]string{
	OrderConfirmation: "Order_Confirmation",
	UserVerification: "User_Verification",
	PassReset: "Pass_Reset",
}

const (
	MailEndpoint string = "/v3/mail/send"
	MailHost 	 string = "https://api.sendgrid.com"
)

func(sg *SendgridManager) CreateMail(mailReq *Mail) []byte {
	m := mail.NewV3Mail()

	from := mail.NewEmail("Flommerce", mailReq.from)
	m.SetFrom(from)

	switch mailReq.mailTyp {
		case UserVerification: {
			m.SetTemplateID(sg.TUserVerification)
		}
		case PassReset: {
			m.SetTemplateID(sg.TPassReset)
		}
		case OrderConfirmation: {
			m.SetTemplateID(sg.TOrderConfirmation)
		}
	default:
		logger_system.Log.Error("Unknown Mail type", zap.Int("type", int(mailReq.mailTyp)))
		return nil
	}

	p := mail.NewPersonalization()

	tos := make([]*mail.Email, 0)
	for _, to := range mailReq.to {
		tos = append(tos, mail.NewEmail("user", to))
	}

	p.AddTos(tos...)
	p.SetDynamicTemplateData("domainName", mailReq.data.DomainName)
	p.SetDynamicTemplateData("firstName", mailReq.data.FirstName)
	p.SetDynamicTemplateData("username", mailReq.data.Username)
	p.SetDynamicTemplateData("verificationUrl", mailReq.data.TokenUrl)
	p.SetDynamicTemplateData("supportEmail", mailReq.data.SupportEmail)
	p.SetDynamicTemplateData("currentYear", time.Now().Year())

	m.AddPersonalizations(p)
	return mail.GetRequestBody(m)
}

func(sg *SendgridManager) SendMail(mailReq *Mail) error {
	request := sendgrid.GetRequest(sg.SendgridApiSecret, MailEndpoint, MailHost)
	request.Method = "POST"

	body := sg.CreateMail(mailReq)
	request.Body = body
	response, err := sendgrid.API(request)
	if err != nil {
		logger_system.Log.Error("Unable to send email", zap.Error(err))
		return err
	}

	logger_system.Log.Info("Mail sent successfully", zap.String("Status code", strconv.Itoa(response.StatusCode)))

	return nil
}

func(sg *SendgridManager) NewMail(params *MailServiceParams) *Mail {
	return &Mail{
		from: params.From,
		to: params.To,
		subject: params.Subject,				
		mailTyp: params.MailType,
		data: params.MailData,
	}
}

func(sg *SendgridManager) Validate() error {
	switch {
	case sg.TUserVerification == "" : 
		return fmt.Errorf("Missing User Verification template")		
	case sg.TPassReset == "": 
		return fmt.Errorf("Missing Password Reset template")		
	case sg.TOrderConfirmation == "": 
		return fmt.Errorf("Missing Order template")		
	case sg.TestEmail == "": 
		return fmt.Errorf("Missing test email")		
	case sg.DomainEmail == "": 
		return fmt.Errorf("Missing domain email")
	case sg.CustomerSupportEmail == "": 
		return fmt.Errorf("Missing Customer Support Email")
	default:
		return nil
	}
} 