package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrInvalidToken       = errors.New("invalid token")
    ErrExpiredToken       = errors.New("token has expired")
    ErrEmailInUse         = errors.New("email already in use")
)

type AuthClaims struct {
    UserID          string      `json:"userId"`
    UserType        string      `json:"userType"`
    UserVerified    bool        `json:"userVerified"`
    jwt.RegisteredClaims
}

type UserVerificationClaims struct {
    UserID          string      `json:"userId"`
    MailType        string      `json:"mailType"`
    jwt.RegisteredClaims
}; 

type SudoClaims struct {
    UserID          string       `json:"userId"`
    UserType        string       `json:"userType"`
    jwt.RegisteredClaims
}

func GenerateAccessToken(
    userID, 
    userType string, 
    verifiedStatus bool, 
    secret []byte,
) (string, error) {
    if userID == "" || userType == "" {
        return "", fmt.Errorf("Missing required identifier")
    }

    claims := AuthClaims{
        UserID: userID,
        UserType: userType,
        UserVerified: verifiedStatus,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt: jwt.NewNumericDate(time.Now()),            
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
    	return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(
    userID, 
    userType string, 
    secret []byte,
) (string, time.Time, error) {
    expiresAt := time.Now().Add(3 * 24 * time.Hour)
    claims := &AuthClaims{
        UserID: userID,
        UserType: userType,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt: jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    signed ,err := token.SignedString(secret)
    return signed, expiresAt, err
}

func GenerateSudoToken(
    user_id, 
    user_type string, 
    secret []byte,
)(string, error) {
    if user_id == "" && user_type == "" {
        return "", fmt.Errorf("Missing required identifier")
    }

    claims := SudoClaims{
        UserID: user_id,
        UserType: user_type,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
            IssuedAt: jwt.NewNumericDate(time.Now()),
        },        
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
    	return "", err
    }
        
	return tokenString, nil
}

func GenerateUserVerificationToken(
    user_id string, 
    secret []byte,
) (string, time.Time, error) {    
    expiresAt := time.Now().Add(24 *time.Hour)

    claims := UserVerificationClaims{
        UserID: user_id,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt: jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)
    signed, err := token.SignedString(secret)

    return signed, expiresAt, err
}

func VerifyAccessToken(tokenString string, secret []byte) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		} 
			return secret, nil 
	})
	
    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrExpiredToken
        }

        return nil, ErrInvalidToken
    }

	claims, ok := token.Claims.(*AuthClaims)
    if !ok || !token.Valid {
        return nil, ErrInvalidToken
    }

	return claims, nil
}

func VerifySudoToken(tokenString string, secret []byte) (*SudoClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &SudoClaims{}, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidToken
        }
        return secret, nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrExpiredToken
        }

        return nil, ErrInvalidToken
    }

    claims, ok := token.Claims.(*SudoClaims) 
    if !ok || !token.Valid {
        return nil, ErrInvalidToken
    }

    return claims, nil
}

func VerifyUserVerificationToken(tokenString string, secret []byte)(*UserVerificationClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &UserVerificationClaims{}, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidToken
        }
        return secret, nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired){
            return nil, ErrExpiredToken
        }

        return nil, ErrInvalidToken
    }

    claims, ok := token.Claims.(*UserVerificationClaims)
    if !ok || !token.Valid {
        return nil, ErrInvalidToken
    }

    return claims, nil
}