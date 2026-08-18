package user_type

import (
	"backend/src/models"
	"time"
)


type UserRegisterRequest struct {
    FirstName    	string 	`json:"firstName"   binding:"required"`
    LastName     	string 	`json:"lastName"    binding:"required"`
    Username     	string 	`json:"username"    binding:"required"`

    Email       	string 	`json:"email"       binding:"required,email"`
	PhoneNumber		string 	`json:"phoneNumber" binding:"required"`
    Password     	string 	`json:"password"    binding:"required,min=8"`
    UserLocale 	 	string 	`json:"userLocale"  binding:"required"`
	UserCountry  	string 	`json:"userCountry" binding:"required"`
	UserAddress  	string 	`json:"userAddress" binding:"required"`
	
	IsAgreed	 	bool 	`json:"isAgreed" binding:"required"`
	EmailConsent  	bool	`json:"emailConsent" binding:"omitempty"`
	SmsConsent	  	bool	`json:"smsConsent" binding:"omitempty"`	
	ConsentSource 	string 	`json:"consentSrc" binding:"omitempty"`	
}

type SearchUserRequest struct {
	Query		*string 	`form:"q"`
	Username	*string		`form:"username"`
	SortBy		*string		`form:"sortBy"`
	SortOrder 	*string		`form:"sortOrder"`
	Cursor		*string		`form:"cursor"`
	Limit 		int 		`form:"limit"`
}

type UpdateUserRequest struct {
    NewFirstname 	*string `json:"newFirstname" binding:"omitempty"`
    NewLastname  	*string `json:"newLastname"  binding:"omitempty"`
	NewUsername  	*string `json:"newUsername"  binding:"omitempty,min=3"`

	NewPhonenumber	*string	`json:"newPhonenumber" binding:"omitempty"`
	NewEmail     	*string `json:"newEmail"     binding:"omitempty,email"`
    NewPassword  	*string `json:"newPassword"  binding:"omitempty,min=8"`

    NewLocale  		*string `json:"newLocale" binding:"omitempty"`
	NewCountry		*string	`json:"newCountry" binding:"omitempty"`
	NewAddress		*string `json:"newAddress" binding:"omitempty"`

	NewEmailConsent *bool 	`json:"newEmailConsent" binding:"omitempty"`
	NewSmsConsent	*bool	`json:"newSmsConsent" binding:"omitempty"`
}

type userResponse struct {
	UserID       	string    	`json:"userId"`
    FirstName    	string    	`json:"firstName"`
    LastName    	string    	`json:"lastName"`
    Username   	  	string    	`json:"username"`
	PhoneNumber		string	  	`json:"phoneNumber"`
    Email        	string    	`json:"email"`
    UserType     	string    	`json:"userType"`
    UserLocale 	 	string    	`json:"userLocale"`
	UserCountry  	string	  	`json:"userCountry"`
	UserAddress  	string    	`json:"userAddress"`
	IsAgreed		bool 	  	`json:"isAgreed"`
    IsVerified   	bool      	`json:"isVerified"` 
    CreatedAt    	time.Time 	`json:"createdAt"`
	UpdatedAt 	 	*time.Time 	`json:"updatedAt"`
	ConsentUpdated 	*time.Time 	`json:"consentUpdated"`
}

func CreateUserResponse(u *models.User) userResponse {
    return userResponse{
        UserID:       u.UserID,
        FirstName:    u.FirstName,
        LastName:     u.LastName,
        Username:     u.Username,
        Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
        UserType:     u.UserType,
        UserLocale:   u.Locale,
		UserCountry:  u.Country,
		UserAddress:  u.Address,
        IsVerified:   u.IsVerified,
		IsAgreed: 	  u.IsAgree,
        CreatedAt:    u.CreatedAt,
		ConsentUpdated: u.Consent_Updated,
		UpdatedAt: u.Updatedat,
    }
}