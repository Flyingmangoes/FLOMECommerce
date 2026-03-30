package middlewares

import (
    "log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
    Code     int    `json:"code"`
    Message  string `json:"message"`
}

func (e *AppError) Error() string {
    return e.Message
}

// Helper constructors
func ErrBadRequest(msg string) *AppError {
    return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func ErrConflict(msg string) *AppError {
    return &AppError{Code: http.StatusConflict, Message: msg}
}

func ErrInternal(msg string) *AppError {
    return &AppError{Code: http.StatusInternalServerError, Message: msg}
}

func ErrUnauthorized(msg string) *AppError {
    return &AppError{Code: http.StatusUnauthorized, Message: msg}
}

func ErrNotFound(msg string) *AppError {
    return &AppError{Code: http.StatusNotFound, Message: msg}
}

func JSONAppErrorReporter() gin.HandlerFunc {
    return jsonAppErrorReporterT(gin.ErrorTypeAny)
}

//
// Middleware Error Handler in server package
//

func jsonAppErrorReporterT(errType gin.ErrorType) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        detectedErrors := c.Errors.ByType(errType)

        if len(detectedErrors) > 0 {
            log.Println("Handle APP error")

            err := detectedErrors[0].Err
            var parsedError *AppError
            switch e := err.(type) {    
            case *AppError:
                parsedError = e
            default:
                parsedError = &AppError{ 
                  Code: http.StatusInternalServerError,
                  Message: "Internal Server Error",
                }
            }
            // Put the error into response
            c.IndentedJSON(parsedError.Code, parsedError)
            c.Abort()
            // or c.AbortWithStatusJSON(parsedError.Code, parsedError)
            return
        }

    }
}