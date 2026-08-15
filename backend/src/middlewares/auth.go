package middlewares

import (
	terror "backend/src/error"
	"backend/src/services/auth"
	jwt_service "backend/src/utils/JWT"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddlewares(secret string) gin.HandlerFunc {
	return func (c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(terror.ErrUnauthorized("Missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2) 
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(terror.ErrUnauthorized("Invalid authorization format"))
			c.Abort()
			return 
		}

		claims, err := jwt_service.VerifyAccessToken(parts[1], []byte(secret))
		if err != nil {
			c.Error(terror.ErrUnauthorized(err.Error()))
            c.Abort()
            return
		}

		c.Set("userVerified", claims.UserVerified)
		c.Set("userId", claims.UserID)
        c.Set("userType", claims.UserType)
        c.Next()
	}
}        

func AuthorizationMiddleware(action auth.Action) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("userId")
		userType := c.GetString("userType")

		ok, err := auth.VerifyAuthorization(userId, auth.AccountType(userType), action)
		if err != nil  || !ok {
			c.Error(terror.ErrUnauthorized("Invalid user"))
			c.Abort()
			return
		}

		c.Next()
	}
}

type LoginAttempt struct {
	Attempts int
	LockedAt *time.Time
}

type LoginPrison struct {
	mu sync.Mutex
	records map[string]*LoginAttempt
	Max int
	Lockout time.Duration
}

func NewLoginPrison(max int, lockout time.Duration) *LoginPrison {
	lp := &LoginPrison{
		records: make(map[string]*LoginAttempt),
		Max: max,
		Lockout: lockout,
	}
	go lp.cleanup()
	return lp
}

func (lp *LoginPrison) RecordFailure(key string) {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	record, exists := lp.records[key]
	if !exists {
		record = &LoginAttempt{}
		lp.records[key] = record
	}

	record.Attempts++
	if record.Attempts >= lp.Max {
		now := time.Now()
		record.LockedAt = &now
	}
}

func (lp *LoginPrison) IsLocked(key string) (bool, time.Duration) {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	record, exists := lp.records[key]
	if !exists {
		return false, 0
	}

	if record.LockedAt == nil {
		return false, 0
	}

	remaining := lp.Lockout - time.Since(*record.LockedAt)
	if remaining <= 0 {
		delete(lp.records, key)
		return false, 0
	}

	return true, remaining
}

func (lp *LoginPrison) Release(key string) {
	lp.mu.Lock()
	defer lp.mu.Unlock()
	delete(lp.records, key)
}

func (lp *LoginPrison) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		lp.mu.Lock()
		for key, record := range lp.records {
			if record.LockedAt != nil && time.Since(*record.LockedAt) > lp.Lockout {
				delete(lp.records, key)
			}
		}
		lp.mu.Unlock()
	}
}