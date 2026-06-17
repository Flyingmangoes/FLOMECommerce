package repository

import (
	"backend/src/models"
	"context"
	"database/sql"
	"fmt"
)

type UserStoreInterface interface {
	CreateUser(ctx context.Context, params *UserProfileParams) (*models.User, error)
	UpdateUser(ctx context.Context, params *UserProfileParams) (*models.User, error)
	DeleteUser(ctx context.Context, params *UserProfileParams) (error)
	LoginUser(ctx context.Context, params *UserProfileParams) (*models.User, error)
	ResetPassword(ctx context.Context, params *UserProfileParams)(error)

	FetchUserByID(ctx context.Context, user_id *string) (*models.User, error)
	FetchPassword(ctx context.Context, params *UserProfileParams) (*string, error)

	VerifyUser(ctx context.Context, verified_id string) (*models.User, error)
	SearchUser(ctx context.Context, params *UserSearchParams) ([]models.User, error)
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
		`INSERT INTO mkt_ecommerce.mkt_users (firstname, lastname, username, 
			email, phone_number, password_hash, user_locale, user_country, user_address, 
			user_type, is_agree, email_consent, sms_consent, consent_src)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING user_id, firstname, lastname, 
		username, email, phone_number, 
		user_locale, user_country, user_address, 
		user_type, is_verified, is_agree, 
		email_consent, sms_consent, consent_src, 
		created_at`,
		params.FirstName, params.LastName, params.Username, 
		params.Email, params.PhoneNumber, 
		params.HashedPassword, params.Locale, params.Country, 
		params.Address, params.UserType, params.IsAgree,
		params.EmailConsent, params.SmsConsent, params.ConsentSource,
	).Scan(&user.UserID, 
		&user.FirstName, &user.LastName, &user.Username, 
		&user.Email, &user.PhoneNumber, &user.Locale,
		&user.Country, &user.Address, &user.UserType, &user.IsAgree, 
		&user.IsAgree, &user.EmailConsent, &user.SmsConsent, 
		&user.Consent_src, &user.CreatedAt,
	)

	if err != nil {
		return nil, err 
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStore)UpdateUser(ctx context.Context, params *UserProfileParams) (*models.User, error) {
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
		RETURNING user_id, firstname, lastname, phone_number, email, password_hash, 
		user_locale, user_country, user_address, email_consent, sms_consent, 
		consent_updated_at, updated_at`,
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

	if err != nil {
		return nil, err 
	}

	return user, nil
}

func (us *UserStore)DeleteUser(ctx context.Context, params *UserProfileParams) error {
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

func (us *UserStore)LoginUser(ctx context.Context, params *UserProfileParams) (*models.User, error) {
	user := &models.User{}
	var err error

	hasEmail := params.Email != nil && *params.Email != ""
	hasUsername := params.Username != nil && *params.Username != ""

	if !hasEmail && !hasUsername {
		return nil, fmt.Errorf("email or username is required")
	}

	const baseSelect = `
        SELECT user_id, firstname, lastname,
            username, email, phone_number, password_hash,
            user_locale, user_country, user_address,
            user_type, is_verified, is_agree,
            email_consent, sms_consent, consent_src,
            created_at, updated_at, consent_updated_at
        FROM mkt_ecommerce.mkt_users`

	if hasEmail {
		err = us.db.QueryRowContext(ctx,
		baseSelect + ` WHERE email = $1`,
		*params.Email,
		).Scan(
			&user.UserID, &user.FirstName, &user.LastName,
            &user.Username, &user.Email, &user.PhoneNumber, &user.PasswordHash,
            &user.Locale, &user.Country, &user.Address,
            &user.UserType, &user.IsVerified, &user.IsAgree,
            &user.EmailConsent, &user.SmsConsent, &user.Consent_src,
            &user.CreatedAt, &user.Updatedat, &user.Consent_Updated,
		)
	} else {
		err = us.db.QueryRowContext(ctx,
			baseSelect + ` WHERE username = $1`,
			*params.Username,
		).Scan(
			&user.UserID, &user.FirstName, &user.LastName,
            &user.Username, &user.Email, &user.PhoneNumber, &user.PasswordHash,
            &user.Locale, &user.Country, &user.Address,
            &user.UserType, &user.IsVerified, &user.IsAgree,
            &user.EmailConsent, &user.SmsConsent, &user.Consent_src,
            &user.CreatedAt, &user.Updatedat, &user.Consent_Updated,
		)
	}

	if err != nil {
		return nil, err 
	}

	return user, nil
}

func (us *UserStore) FetchUserByID(ctx context.Context, user_id *string) (*models.User, error) {
	user := &models.User{}

	err := us.db.QueryRowContext(ctx,
		`SELECT user_id, firstname, lastname, username, email, phone_number, 
			user_locale, user_country, user_address, user_type, 
			is_verified, is_agree, email_consent, sms_consent, 
			consent_src, created_at, 
			updated_at, consent_updated_at 
		FROM mkt_ecommerce.mkt_users 
		WHERE user_id = $1`,
		user_id,
	).Scan(&user.UserID, 
		&user.FirstName, &user.LastName, &user.Username, 
		&user.Email, &user.PhoneNumber, &user.Locale, 
		&user.Country, &user.Address, &user.UserType, 
		&user.IsVerified, &user.IsAgree, &user.EmailConsent, &user.SmsConsent, 
		&user.Consent_src, &user.CreatedAt, &user.Updatedat, 
		&user.Consent_Updated,
	)
	if err != nil {
		return nil, err 
	}

	return user, nil
}

func (us *UserStore) FetchPassword(ctx context.Context, params *UserProfileParams) (*string, error) {
    user := &models.User{}
	var err error

	switch {
		case params.Email != nil:
			err = us.db.QueryRowContext(ctx,
				`SELECT password_hash from mkt_ecommerce.mkt_users
				WHERE email = $1`,	
				params.Email,
			).Scan(&user.PasswordHash)
		
		case params.UserId != nil:
			err = us.db.QueryRowContext(ctx,
				`SELECT password_hash FROM mkt_ecommerce.mkt_users 
				WHERE user_id = $1`,
				params.UserId,
			).Scan(&user.PasswordHash)

		case params.Username != nil:
			err = us.db.QueryRowContext(ctx,
				`SELECT password_hash FROM mkt_ecommerce.mkt_users
				WHERE username = $1`,
				params.Username,
			).Scan(&user.PasswordHash)

		default:
			return nil, fmt.Errorf("at least one identifier required")
	}

	if err != nil {
		return nil, err 
	}

    return &user.PasswordHash, nil
}		

func (us *UserStore) VerifyUser(ctx context.Context, verified_id string) (*models.User, error) {
	user := &models.User{}

	err := us.db.QueryRowContext(ctx,
		`UPDATE mkt_ecommerce.mkt_users SET
			is_verified = true,
			updated_at = NOW()
		WHERE user_id = $1
		RETURNING user_id, email`,
		verified_id,
	).Scan(&user.UserID, &user.Email)

	if err != nil { return nil, err }

	return user, nil
}

func (us *UserStore) SearchUser(ctx context.Context, params *UserSearchParams) ([]models.User, error) {
	params.Normalize()

	createdAt, id := params.CursorValues()
	
	rows, err := us.db.QueryContext(ctx,
		`SELECT user_id, username, user_country, created_at,
		FROM mkt_ecommerce.mkt_users
		WHERE
			($1::timestamptz IS NULL OR (created_at, user_id) < ($1, $2::uuid)) AND
			($3::varchar IS NULL OR username ILIKE '%' || $3 || '%') AND
		ORDER BY created_at DESC, user_id DESC
		LIMIT $4`,
		createdAt, id,
		params.Username, params.Limit,
	)

	if err != nil { return nil, err }
	
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		u := models.User{}
		if err := rows.Scan(
			&u.UserID, &u.Username,
			&u.Country, &u.CreatedAt, 
		); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, rows.Err()
}

func(us *UserStore) ResetPassword(ctx context.Context, params *UserProfileParams) error {
	var results sql.Result; var err error

	const query string = `
		UPDATE mkt_ecommerce.mkt_users SET
			password_hash = COALESCE($1, password_hash),
			updated_at = NOW()
	`

	switch {
		case params.Email != nil:{
			results, err = us.db.ExecContext(
				ctx, 
				query + `WHERE email = $1`, 
				params.Email,
			)
		}
		case params.Username != nil:{
			results, err = us.db.ExecContext(
				ctx,
				query + `WHERE username = $1`,
				params.Username,
			)
		}
		default: {
			return fmt.Errorf("at least one identifier is required")
		}
	}

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