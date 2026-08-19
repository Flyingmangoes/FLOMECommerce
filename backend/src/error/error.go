package error_service

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

//
// Middleware Error Handler in server package
//

func JSONAppErrorReporter() gin.HandlerFunc {
    return jsonAppErrorReporterT(gin.ErrorTypeAny)
}

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