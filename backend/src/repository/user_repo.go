package repository

import (
	"backend/src/models"
	repo_type "backend/src/repository/types"
	"context"
	"database/sql"
	"fmt"
)

const (
	UserVerified 	= iota + 1
	UserUnverified 	
	UserMerchant 	
	UserAdmin	
)

var user_type_list map[int]string = map[int]string{
	UserVerified: 	"VERIFIED_USER",
	UserUnverified: "UNVERIFIED_USER",
	UserMerchant: 	"MERCHANT",
	UserAdmin:		"ADMIN",
}

type UserStoreInterface interface {
	Create(ctx context.Context, params *repo_type.UserProfileParams) (*models.User, error)
	Update(ctx context.Context, params *repo_type.UserProfileParams) (*models.User, error)
	Delete(ctx context.Context, params *repo_type.UserProfileParams) (error)
	Login(ctx context.Context, params *repo_type.UserProfileParams) (*models.User, error)
	ResetPassword(ctx context.Context, user_id string)(error)

	Fetch(ctx context.Context, user_id string) (*models.User, error)
	FetchPassword(ctx context.Context, user_id string) (*string, error)

	VerifyUser(ctx context.Context, verified_user_id string) (*models.User, error)
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (us *UserStore)Create(ctx context.Context, params *repo_type.UserProfileParams) (*models.User, error) {
	const query string = `
		INSERT INTO mkt_ecommerce.mkt_users (firstname, lastname, username, 
				email, phone_number, password_hash, 
				user_locale, user_country, user_address, 
				user_type, is_agree, email_consent, 
				sms_consent, consent_src)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			RETURNING user_id, firstname, lastname, 
			username, email, phone_number, user_type,
			user_locale, user_country, user_address, 
			is_verified, is_agree, 
			email_consent, sms_consent, consent_src, 
			created_at
	`
	user := &models.User{}

	tx, err := us.db.BeginTx(ctx, nil)
	if err != nil { return nil, err } 

	defer tx.Rollback()

	if params.UserType == nil {
		return nil, fmt.Errorf("User Type missing")
	}

	if *params.UserType == UserAdmin {
		err = tx.QueryRowContext(
			ctx, 
			query,
			params.FirstName, 	params.LastName, 
			params.Username,  	params.Email, 
			params.PhoneNumber, params.HashedPassword, 
			params.Locale, 		params.Country, 
			params.Address, 	user_type_list[UserAdmin], 
			params.IsAgree, 	params.EmailConsent, 
			params.SmsConsent, 	params.ConsentSource,
		).Scan(
			&user.UserID, 
			&user.FirstName, 	&user.LastName, 
			&user.Username,  	&user.Email, 
			&user.PhoneNumber, 	&user.UserType,
			&user.Locale, 		&user.Country, 		
			&user.Address,  	&user.IsAgree, 		
			&user.IsAgree, 
			&user.EmailConsent, &user.SmsConsent, 
			&user.Consent_src, 	&user.CreatedAt,
		)
	} else {
		err = tx.QueryRowContext(
			ctx, 
			query,
			params.FirstName, 	params.LastName, 
			params.Username,  	params.Email, 
			params.PhoneNumber, params.HashedPassword, 
			params.Locale, 		params.Country, 
			params.Address, 	user_type_list[UserUnverified], 
			params.IsAgree, 	params.EmailConsent, 
			params.SmsConsent, 	params.ConsentSource,
		).Scan(
			&user.UserID, 
			&user.FirstName, 	&user.LastName, 
			&user.Username,  	&user.Email, 
			&user.PhoneNumber, 	&user.UserType,
			&user.Locale, 		&user.Country, 		
			&user.Address,  	&user.IsAgree, 		
			&user.IsAgree, 
			&user.EmailConsent, &user.SmsConsent, 
			&user.Consent_src, 	&user.CreatedAt,
		)
	}

	if err != nil {
		return nil, err 
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStore)Update(ctx context.Context, params *repo_type.UserProfileParams) (*models.User, error) {
	user:= &models.User{}
		
	err := us.db.QueryRowContext(ctx, 
		`UPDATE mkt_ecommerce.mkt_users SET 
			firstname			= COALESCE ($1, firstname),
			lastname			= COALESCE ($2, lastname),
			username			= COALESCE ($3, username),
			phone_number		= COALESCE ($4, phone_number),
			email				= COALESCE ($5, email),
			password_hash 		= COALESCE ($6, password_hash),		
			user_locale  		= COALESCE ($7, user_locale),
			user_country		= COALESCE ($8, user_country),
			user_address 		= COALESCE ($9, user_address),
			email_consent		= COALESCE ($10, email_consent),
			sms_consent			= COALESCE ($11, sms_consent),
			consent_updated_at 	= NOW(),
			updated_at			= NOW()
		WHERE user_id = $12
		RETURNING user_id, firstname, lastname, 
		phone_number, email, password_hash, 
		user_locale, user_country, user_address, 
		email_consent, sms_consent, consent_updated_at, 
		updated_at`,
		params.FirstName, params.LastName, params.Username, 
		params.PhoneNumber, params.Email, params.NewPasswordHashed, 
		params.Locale, params.Country, params.Address, 
		params.EmailConsent, params.SmsConsent, params.UserId,
	).Scan(
		&user.UserID, 		&user.FirstName, 
		&user.LastName, 	&user.PhoneNumber, 
		&user.Email, 		&user.PasswordHash, 
		&user.Locale, 		&user.Country,
		&user.Address, 		&user.EmailConsent, 
		&user.SmsConsent, 	&user.Consent_Updated, 
		&user.Updatedat,
	)

	if err != nil {
		return nil, err 
	}

	return user, nil
}

func (us *UserStore)Delete(ctx context.Context, params *repo_type.UserProfileParams) error {
	result, err := us.db.ExecContext(ctx, 
		`DELETE FROM mkt_ecommerce.mkt_users
		WHERE user_id = $1`,
		params.UserId, 
	)

	if err != nil {
		return err 
	}

	rows, err := result.RowsAffected()
	if err != nil { return err }
	if rows == 0 {
		return fmt.Errorf("User not found or credentials didn't match")
	}

	return nil
}

func (us *UserStore)Login(ctx context.Context, params *repo_type.UserProfileParams) (*models.User, error) {
	user := &models.User{}
	var err error

	const baseSelect = `
        SELECT user_id, firstname, lastname, 
			username, email, phone_number, user_type,
			user_locale, user_country, user_address, 
			is_verified, is_agree,
			email_consent, sms_consent, consent_src, 
			created_at, updated_at, consent_updated_at,
        FROM mkt_ecommerce.mkt_users`

	switch {
	case params.Email != nil && params.Username == nil: 
		err = us.db.QueryRowContext(ctx,
		baseSelect + ` WHERE email = $1`,
		*params.Email,
		).Scan(
			&user.UserID, 		&user.FirstName, 
			&user.LastName, 	&user.Username, 
			&user.Email, 		&user.PhoneNumber,
			&user.UserType, 
			&user.Locale, 		&user.Country,
			&user.Address, 		&user.IsVerified, 
			&user.IsAgree,		&user.EmailConsent,
			&user.SmsConsent, 	&user.CreatedAt, 
			&user.Updatedat, 	&user.Consent_Updated,
		)
	case params.Username != nil && params.Email == nil: 
		err = us.db.QueryRowContext(ctx,
			baseSelect + ` WHERE username = $1`,
			*params.Username,
		).Scan(
			&user.UserID, 		&user.FirstName, 
			&user.LastName, 	&user.Username, 
			&user.Email, 		&user.PhoneNumber,
			&user.UserType, 
			&user.Locale, 		&user.Country,
			&user.Address, 		&user.IsVerified, 
			&user.IsAgree,		&user.EmailConsent,
			&user.SmsConsent, 	&user.CreatedAt, 
			&user.Updatedat, 	&user.Consent_Updated,
		)
	default:
		return nil, fmt.Errorf("at least one identifier required")
	}

	if err != nil {
		return nil, err 
	}

	return user, nil
}

func (us *UserStore) Fetch(ctx context.Context, user_id string) (*models.User, error) {
	user := &models.User{}

	if user_id == "" {
		return nil, fmt.Errorf("user_id required")
	}

	err := us.db.QueryRowContext(ctx,
		`SELECT user_id, firstname, lastname, 
			username, email, phone_number, user_type,
			user_locale, user_country, user_address, 
			is_verified, is_agree,
			email_consent, sms_consent, consent_src, 
			created_at, updated_at, consent_updated_at,
		FROM mkt_ecommerce.mkt_users 
		WHERE user_id = $1`,
		user_id,
	).Scan(
		&user.UserID, 		&user.FirstName, 
		&user.LastName, 	&user.Username, 
		&user.Email, 		&user.PhoneNumber,
		&user.UserType, 
		&user.Locale, 		&user.Country,
		&user.Address, 		&user.IsVerified, 
		&user.IsAgree,		&user.EmailConsent,
		&user.SmsConsent, 	&user.CreatedAt, 
		&user.Updatedat, 	&user.Consent_Updated,
	)
	if err != nil {
		return nil, err 
	}

	return user, nil
}

func (us *UserStore) FetchPassword(ctx context.Context, user_id string) (*string, error) {
    user := &models.User{}

	if user_id == "" {
		return nil, fmt.Errorf("user_id required")
	}
		
	err := us.db.QueryRowContext(ctx,
		`SELECT password_hash FROM mkt_ecommerce.mkt_users 
		WHERE user_id = $1`,
		user_id,
	).Scan(&user.PasswordHash)

	if err != nil {
		return nil, err 
	}

    return &user.PasswordHash, nil
}		

func (us *UserStore) VerifyUser(ctx context.Context, verified_user_id string) (*models.User, error) {
	user := &models.User{}

	err := us.db.QueryRowContext(ctx,
		`UPDATE mkt_ecommerce.mkt_users SET
			is_verified = true,
			user_type = $1,
			updated_at = NOW()
		WHERE user_id = $2
		RETURNING user_id, email`,
		user_type_list[UserVerified], verified_user_id,
	).Scan(&user.UserID, &user.Email)

	if err != nil { return nil, err }

	return user, nil
}

func(us *UserStore) ResetPassword(ctx context.Context, user_id string) error {
	const query string = `
		UPDATE mkt_ecommerce.mkt_users SET
			password_hash = COALESCE($1, password_hash),
			updated_at = NOW()
	`

	if user_id == "" {
		return fmt.Errorf("user id is required")
	}

	results, err := us.db.ExecContext(ctx, 
		query + `WHERE user_id = $1`,
		user_id,
	)

	if err != nil {
		return err
	}

	rows, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("User not found or credentials didn't match")
	}

	return nil
}