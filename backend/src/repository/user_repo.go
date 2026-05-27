package repository

import (
	"backend/src/models"
	"backend/src/utils"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UserStoreInterface interface {
	CreateUser(ctx context.Context, params *UserProfileParams) (*models.User, error)
	UpdateUser(ctx context.Context, params *UserProfileParams) (*models.User, error)
	DeleteUser(ctx context.Context, params *UserProfileParams) (error)
	LoginByUserEmail(ctx context.Context, email *string) (*models.User, error)

	GetUserByUsername(ctx context.Context, username *string) (*models.User, error)
	GetUserByID(ctx context.Context, user_id *string) (*models.User, error)
	GetPassword(ctx context.Context, params *UserProfileParams) (*models.User, error)
}

type BaseParams struct {
	UserId 			*string
	Email 			*string
	Username 		*string

	Locale			*string
	Country			*string
	Address 		*string
}

type UserProfileParams struct {	
	BaseParams
	PhoneNumber		*string
	HashedPassword  *string

	//For updating (only for the password)
	NewPasswordHashed 	*string

	//Profile Section
	FirstName		*string
	LastName		*string
	UserType 		*string

	//Consent related Section
	EmailConsent	*bool
	SmsConsent		*bool
	ConsentUpdatedAt *time.Time
	ConsentSource	*string

	//Extra Section
	IsAgree			*bool
	IsVerified		*bool	
	CreatedAt		*time.Time
	UpdatedAt		*time.Time
}



type ListUsersFilter struct {
	UserType *string
	utils.PagFilter
}



type UserStore struct {
	db *sql.DB
}



func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}



func (us *UserStore)CreateUser(ctx context.Context, params *UserProfileParams) (*models.User, error) {
	user := &models.User{}

	tx, err := us.db.BeginTx(ctx, nil)
	if err != nil { return nil, err } 

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, 
		`INSERT INTO mkt_users (firstname, lastname, username, email, phone_number, passwordhashed, user_locale, user_country, user_address, user_type, is_agree, email_consent, sms_consent, consent_src)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING user_id, firstname, lastname, username, email, phone_number, user_locale, user_country, user_address, user_type, is_verified, is_agree, email_consent, sms_consent, consent_src, created_at`,
		params.FirstName, params.LastName, params.Username, 
		params.Email, params.PhoneNumber, 
		params.HashedPassword, params.Locale, params.Country, 
		params.Address, params.UserType, params.IsAgree,
	 	params.EmailConsent, params.SmsConsent, params.ConsentSource,
	).Scan(&user.UserID, 
		&user.FirstName, &user.LastName, &user.Username, 
		&user.Email, &user.PhoneNumber, &user.Locale,
		&user.Country, &user.Address, &user.UserType, &user.IsVerified, 
		&user.IsAgree, &user.EmailConsent, &user.SmsConsent, 
		&user.Consent_src, &user.CreatedAt,
	)
	if err != nil { return nil, err }
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStore)UpdateUser(ctx context.Context, params *UserProfileParams) (*models.User, error) {
	user:= &models.User{}
		
	err := us.db.QueryRowContext(ctx, 
		`UPDATE mkt_users SET 
			firstname		= COALESCE ($1, firstname),
			lastname		= COALESCE ($2, lastname),
			username		= COALESCE ($3, username),
			phone_number	= COALESCE ($4, phone_number),
			email			= COALESCE ($5, email),
			passwordhashed 	= COALESCE ($6, passwordhashed),		
			user_locale  	= COALESCE ($7, user_locale),
			user_country	= COALESCE ($8, user_country),
			user_address 	= COALESCE ($9, user_address),
			email_consent	= COALESCE ($10, email_consent),
			sms_consent		= COALESCE ($11, sms_consent),
			consent_updated_at = NOW(),
			updated_at		= NOW()
		WHERE user_id = $12
		RETURNING user_id, firstname, lastname, phone_number, email, passwordhashed, user_locale, user_country, user_address, email_consent, sms_consent, consent_updated_at, updated_at`,
		params.FirstName, params.LastName, params.Username, 
		params.PhoneNumber, params.Email, params.NewPasswordHashed, 
		params.Locale, params.Country, params.Address, 
		params.EmailConsent, params.SmsConsent, params.UserId,
	).Scan(&user.UserID, &user.FirstName, 
		&user.LastName, &user.PhoneNumber, &user.Email, 
		&user.PasswordHash, &user.Locale, &user.Country,
		&user.Address, &user.EmailConsent, &user.SmsConsent, 
		&user.Consent_Updated, &user.Updatedat,
	)
	if err != nil { return nil, err }

	return user, nil
}



func (us *UserStore)DeleteUser(ctx context.Context, params *UserProfileParams) error {
	result, err := us.db.ExecContext(ctx, 
		`DELETE FROM mkt_users
		WHERE user_id = $1 AND email = $2`,
		params.UserId, params.Email,
	)
	if err != nil { return err }

	rows, err := result.RowsAffected()
	if err != nil { return err }
	if rows == 0 {
		return fmt.Errorf("user not found or credentials do not match")
	}

	return nil
}



func (us *UserStore)LoginByUserEmail(ctx context.Context, email *string) (*models.User, error) {
	user := &models.User{}

	err := us.db.QueryRowContext(ctx,
		`SELECT user_id, firstname, lastname, username, email, phone_number, passwordhashed, user_locale, user_country, user_address, user_type, is_verified, is_agree, email_consent, sms_consent, consent_src, created_at, updated_at, consent_updated_at FROM mkt_users
		WHERE email = $1`,
		email,
	).Scan(&user.UserID, 
		&user.FirstName, &user.LastName, &user.Username, 
		&user.Email, &user.PhoneNumber, &user.PasswordHash, 
		&user.Locale, &user.Country, &user.Address, &user.UserType, 
		&user.IsVerified, &user.IsAgree, &user.EmailConsent, &user.SmsConsent, 
		&user.Consent_src, &user.CreatedAt, &user.Updatedat, 
		&user.Consent_Updated,
	)

	if err != nil { return nil, err }

	return user, nil
}



func (us *UserStore) GetUserByUsername(ctx context.Context, username *string) (*models.User, error) {
	user := &models.User{}

	err := us.db.QueryRowContext(ctx,
		`SELECT user_id, firstname, lastname, username, email, phone_number, passwordhashed, user_locale, user_country, user_address, user_type, is_verified, is_agree, email_consent, sms_consent, consent_src, created_at, updated_at, consent_updated_at FROM mkt_users 
		WHERE username = $1`,
		username,
	).Scan(&user.UserID, 
		&user.FirstName, &user.LastName, &user.Username, 
		&user.Email, &user.PhoneNumber, &user.PasswordHash, 
		&user.Locale, &user.Country, &user.Address, &user.UserType, 
		&user.IsVerified, &user.IsAgree, &user.EmailConsent, &user.SmsConsent, 
		&user.Consent_src, &user.CreatedAt, &user.Updatedat, 
		&user.Consent_Updated,
	)
	if err != nil { return nil, err }

	return user, nil
}

func (us *UserStore) GetUserByID(ctx context.Context, user_id *string) (*models.User, error) {
	user := &models.User{}

	err := us.db.QueryRowContext(ctx,
		`SELECT user_id, firstname, lastname, username, email, phone_number, passwordhashed, user_locale, user_country, user_address, user_type, is_verified, is_agree, email_consent, sms_consent, consent_src, created_at, updated_at, consent_updated_at FROM mkt_users 
		WHERE user_id = $1`,
		user_id,
	).Scan(&user.UserID, 
		&user.FirstName, &user.LastName, &user.Username, 
		&user.Email, &user.PhoneNumber, &user.PasswordHash, 
		&user.Locale, &user.Country, &user.Address, &user.UserType, 
		&user.IsVerified, &user.IsAgree, &user.EmailConsent, &user.SmsConsent, 
		&user.Consent_src, &user.CreatedAt, &user.Updatedat, 
		&user.Consent_Updated,
	)
	if err != nil { return nil, err }

	return user, nil
}

func (us *UserStore) GetPassword(ctx context.Context, params *UserProfileParams) (*models.User, error) {
    user := &models.User{}

	var err error

	switch {
		case params.Email != nil:
			err = us.db.QueryRowContext(ctx,
				`SELECT passwordhashed from mkt_users
				WHERE email = $1`,	
				params.Email,
			).Scan(&user.PasswordHash)
		
		case params.UserId != nil:
			err = us.db.QueryRowContext(ctx,
				`SELECT passwordhashed FROM mkt_users 
				WHERE user_id = $1`,
				params.UserId,
			).Scan(&user.PasswordHash)

		case params.Username != nil:
			err = us.db.QueryRowContext(ctx,
				`SELECT passwordhashed FROM mkt_users
				WHERE username = $1`,
				params.Username,
			).Scan(&user.PasswordHash)

		default:
			return nil, fmt.Errorf("at least one identifier must be provided")
	}
	if err != nil { return nil, err }

    return user, nil
}		