package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrInvalidToken       = errors.New("invalid token")
    ErrExpiredToken       = errors.New("token has expired")
    ErrEmailInUse         = errors.New("email already in use")
)

type Claims struct {
    UserID      string `json:"userId"`
    UserType    string `json:"userType"`
    jwt.RegisteredClaims
}

type ConfirmClaims struct {
    UserID string `json:"userId"`
    Action string `json:"action"`
    jwt.RegisteredClaims
}

func GenerateAnyToken(userID, userType, usage string, action *string, secret []byte) (string, error) {
    if usage == "confirm" && action != nil {
        claims := ConfirmClaims{
            UserID: userID,
            Action: *action,
            RegisteredClaims: jwt.RegisteredClaims{
                ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
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

    claims := Claims{
        UserID:   userID,
        UserType: userType,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(userID, userType string, secret []byte) (string, time.Time, error) {
    expiresAt := time.Now().Add(7 * 24 * time.Hour)
    claims := &Claims{
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

func VerifyAccessToken(tokenString string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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

	claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, ErrInvalidToken
    }

	return claims, nil
}

func VerifyConfirmToken(tokenString string, secret []byte) (*ConfirmClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &ConfirmClaims{}, func(t *jwt.Token) (interface{}, error) {
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

    claims, ok := token.Claims.(*ConfirmClaims) 
    if !ok || !token.Valid {
        return nil, ErrInvalidToken
    }

    return claims, nil
}