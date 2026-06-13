package email_services

import (
	"backend/src/repository"
	Logger "backend/src/utils/logger"
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

type SGMailManager struct {
	Users 				repository.UserStoreInterface
	SG_SECRET 				string
	VERIFICATION_SECRET 	string
	TEST_EMAIL 				string
	DOMAIN_EMAIL 			string
	SGTemplate
}	

type SGTemplate struct{
	T_UserVerification 		string
	T_PassReset 			string
	T_OrderConfirmation 	string
}

func NewSGMailManager(sg_secret, v_secret, test_email, domain_email string, template SGTemplate, us repository.UserStoreInterface) *SGMailManager {
	return &SGMailManager{
		SG_SECRET: sg_secret,
		VERIFICATION_SECRET: v_secret,
		Users: us,
		SGTemplate: template,
		TEST_EMAIL: test_email,
		DOMAIN_EMAIL: domain_email,
	}
}

type MailData struct {
	Username 		string		`json:"username"`
	UserEmail 		string		`json:"userEmail"`
	SupportEmail 	string		`json:"supportEmail"`
	Token 			string		`json:"token"`
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
	p.SetDynamicTemplateData("Username", mailReq.data.Username)
	p.SetDynamicTemplateData("User_Email", mailReq.data.UserEmail)
	p.SetDynamicTemplateData("Support_Email", mailReq.data.SupportEmail)
	p.SetDynamicTemplateData("Code", mailReq.data.Token)

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