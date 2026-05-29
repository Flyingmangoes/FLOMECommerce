package email_services

import (
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
	NewMail(params *MailServiceParams) *Mail
}

type MailServiceParams struct {
	from 		string
	to 			[]string
	subject 	string
	mailType 	MailType
	data 		*MailData
}

type SGMailManager struct {
	SG_SECRET 	string
}	

func NewSGMailManager(secret string) *SGMailManager {
	return &SGMailManager{SG_SECRET: secret}
}

type MailData struct {
	Username 		string
	UserEmail 		string
	SupportEmail 	string
	Code 			string
	Expiration 		*time.Duration
}

type Mail struct {
	from 		string
	to 			[]string
	subject 	string
	mailTyp 	MailType
	data 		*MailData
}

const (
	MailConfirmation MailType = iota + 1
	PassReset 
)

const (
	MailEndpoint 	string = "/v3/mail/send"
	MailHost 	 	string = "https://api.sendgrid.com"
)

func(sg *SGMailManager) CreateMail(mailReq *Mail) []byte {
	m := mail.NewV3Mail()

	from := mail.NewEmail("Flommerce", mailReq.from)
	m.SetFrom(from)

	switch mailReq.mailTyp {
		case MailConfirmation: {
			m.SetTemplateID("d-fa0f7c2169aa4c67a6b2f4d9f638a67a")
		}
		case PassReset: {
			m.SetTemplateID("")
		}
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
	p.SetDynamicTemplateData("Link_Expiration", mailReq.data.Expiration)
	p.SetDynamicTemplateData("Code", mailReq.data.Code)

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
		from: params.from,
		to: params.to,
		subject: params.subject,				
		mailTyp: params.mailType,
		data: params.data,
	}
}
