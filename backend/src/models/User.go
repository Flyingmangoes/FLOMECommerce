package models

import "time"

type User struct {
	UserID 			string 		`db:"user_id"`		 
	FirstName 		string 		`db:"firstname"` 
	LastName 		string		`db:"lastname"`
	Username		string		`db:"username"`
	Email 			string 		`db:"email"`
	PhoneNumber		string		`db:"phone_number"`
	PasswordHash 	string 		`db:"passwordhashed"`
	Locale 			string 		`db:"user_locale"`
	Country			string		`db:"user_country"`
	Address			string		`db:"user_address"`
	UserType 		string		`db:"user_type"`
	IsVerified		bool 		`db:"is_verified"`	
	IsAgree			bool		`db:"is_agree"`	
	EmailConsent 	bool		`db:"email_consent"`
	SmsConsent		bool 		`db:"sms_consent"`
	Consent_Updated *time.Time  `db:"consent_updated_at"`
	Consent_src		string		`db:"consent_source"`
	CreatedAt 		time.Time 	`db:"created_at"`
	Updatedat  		*time.Time	`db:"updated_at"`
}	

type VerificationData struct {
	UserID		string 		`db:"user_id"`
	Email		string 		`db:"email"`
	ExpiresAt 	time.Time 	`db:"expires_at"`
	Type 		string		`db:"verification_type"`
}