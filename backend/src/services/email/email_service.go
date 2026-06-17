package email_services

import (
	"backend/src/config"
	"backend/src/repository"
	Logger "backend/src/utils/logger"
	"fmt"
	"strconv"
	"time"
	"net"

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

type SGMailManager struct {
	Users 			repository.UserStoreInterface
	PROD_EMAIL				string
	TEST_EMAIL 				string
	DOMAIN_EMAIL 			string
	SGContext
	SGTemplate
}	

type SGTemplate struct{
	T_UserVerification 		string
	T_PassReset 			string
	T_OrderConfirmation 	string
}

type SGContext struct {
	SERVER_URL 					string
	SG_SECRET 					string
	USER_VERIFICATION_SECRET 	string
}

func NewSGMailManager(cfg *config.ConfigManager, us repository.UserStoreInterface) *SGMailManager {
	return &SGMailManager{
		Users: us,
		TEST_EMAIL: cfg.SENDGRID_CONF.TEST_EMAIL,
		PROD_EMAIL: cfg.SENDGRID_CONF.DOMAIN_EMAIL,
		DOMAIN_EMAIL: cfg.SENDGRID_CONF.DOMAIN_EMAIL,
		SGContext: SGContext{
			SERVER_URL: net.JoinHostPort(cfg.SERV_CONF.HOST, cfg.SERV_CONF.PORT),
			SG_SECRET: cfg.SENDGRID_CONF.SENDGRID_SECRET,
			USER_VERIFICATION_SECRET: cfg.SENDGRID_CONF.VERIFICATION_SECRET,
		},
		SGTemplate: SGTemplate{
			T_UserVerification: cfg.SENDGRID_CONF.TEMPLATE_USER_VERIFICATION,
			T_PassReset: cfg.SENDGRID_CONF.TEMPLATE_PASS_RESET,
			T_OrderConfirmation: cfg.SENDGRID_CONF.TEMPLATE_ORDER_CONFIRMATION,
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

func(sg *SGMailManager) CreateMail(mailReq *Mail) []byte {
	m := mail.NewV3Mail()

	from := mail.NewEmail("Flommerce", mailReq.from)
	m.SetFrom(from)

	switch mailReq.mailTyp {
		case UserVerification: {
			m.SetTemplateID(sg.T_UserVerification)
		}
		case PassReset: {
			m.SetTemplateID(sg.T_PassReset)
		}
		case OrderConfirmation: {
			m.SetTemplateID(sg.T_OrderConfirmation)
		}
	default:
		Logger.Log.Error("Unknown Mail type", zap.Int("type", int(mailReq.mailTyp)))
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

func(sg *SGMailManager) SendMail(mailReq *Mail) error {
	request := sendgrid.GetRequest(sg.SG_SECRET, MailEndpoint, MailHost)
	request.Method = "POST"

	body := sg.CreateMail(mailReq)
	request.Body = body
	response, err := sendgrid.API(request)
	if err != nil {
		Logger.Log.Error("Unable to send email", zap.Error(err))
		return err
	}

	Logger.Log.Info("Mail sent successfully", zap.String("Status code", strconv.Itoa(response.StatusCode)))

	return nil
}

func(sg *SGMailManager) NewMail(params *MailServiceParams) *Mail {
	return &Mail{
		from: params.From,
		to: params.To,
		subject: params.Subject,				
		mailTyp: params.MailType,
		data: params.MailData,
	}
}

func(sg *SGMailManager) Validate() error {
	switch {
		case sg.T_UserVerification == "" : {
			return fmt.Errorf("Missing User Verification template")
		}
		case sg.T_PassReset == "": {
			return fmt.Errorf("Missing Password Reset template")
		}
		case sg.T_OrderConfirmation == "": {
			return fmt.Errorf("Missing Order template")
		}
		case sg.TEST_EMAIL == "": {
			return fmt.Errorf("Missing test email")
		}
		case sg.DOMAIN_EMAIL == "": {
			return fmt.Errorf("Missing domain email")
		}
	}

	return nil
} 