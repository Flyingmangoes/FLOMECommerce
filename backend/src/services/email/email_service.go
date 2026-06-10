package email_services

import (
	"backend/src/repository"
	Logger "backend/src/utils/logger"
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
	SGTemplate
}	

type SGTemplate struct{
	TEMP_EMAILCONFIRMATION 	string
	TEMP_ORDERCONFIRMATION 	string
	TEMP_PASSRESET 		  	string
}

func NewSGMailManager(sg_secret, v_secret string, template SGTemplate, us repository.UserStoreInterface) *SGMailManager {
	return &SGMailManager{
		SG_SECRET: sg_secret,
		VERIFICATION_SECRET: v_secret,
		Users: us,
		SGTemplate: template,
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
	EmailConfirmation MailType = iota + 1
	PassReset 
	OrderConfirmation
)

var typeList = map[MailType]string{
	EmailConfirmation: "Mail_Confirmation",
	PassReset: "Pass_Reset",
}

const (
	MailEndpoint 	string = "/v3/mail/send"
	MailHost 	 	string = "https://api.sendgrid.com"
)

func(sg *SGMailManager) CreateMail(mailReq *Mail) []byte {
	m := mail.NewV3Mail()

	from := mail.NewEmail("Flommerce", mailReq.from)
	m.SetFrom(from)

	switch mailReq.mailTyp {
		case EmailConfirmation: {
			m.SetTemplateID(sg.TEMP_EMAILCONFIRMATION)
		}
		case PassReset: {
			m.SetTemplateID(sg.TEMP_PASSRESET)
		}
		case OrderConfirmation: {
			m.SetTemplateID(sg.TEMP_ORDERCONFIRMATION)
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