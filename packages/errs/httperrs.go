package errs

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stas-bukovskiy/wish-scribe/packages/logger"
	"net/http"
)

type ErrorResponse struct {
	Error ServiceError `json:"error"`
}

type ServiceError struct {
	Kind    string `json:"kind,omitempty"`
	Param   string `json:"param,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewHTTPErrorResponse(ctx *gin.Context, log logger.Logger, err error) {
	if err == nil {
		nilHTTPErrorResponse(ctx, log)
		return
	}

	var e *Error
	if errors.As(err, &e) {
		typicalHTTPErrorResponse(ctx, log, e)
	}

	unknownHTTPErrorResponse(ctx, log, err)
}

func nilHTTPErrorResponse(ctx *gin.Context, log logger.Logger) {
	log.WithContext(ctx).Error("nil error - no response body sent")
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

func typicalHTTPErrorResponse(ctx *gin.Context, log logger.Logger, e *Error) {
	httpStatusCode := httpStatusCode(e.Kind)

	if e.IsEmpty() {
		log.WithContext(ctx).Error("error sent but with empty body")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if httpStatusCode == http.StatusInternalServerError {
		errorLogger(ctx, log, httpStatusCode, e).Error("unexpected error: %+v", e)
	} else {
		errorLogger(ctx, log, httpStatusCode, e).Debug("error response sent to client")
	}

	errorResponse := newErrResponse(e)
	ctx.AbortWithStatusJSON(httpStatusCode, errorResponse)
}

func errorLogger(ctx *gin.Context, log logger.Logger, httpStatusCode int, e *Error) logger.Logger {
	if ops := OpStack(e); len(ops) > 0 {
		return log.WithContext(ctx).With(map[string]interface{}{
			"httpStatusCode": httpStatusCode,
			"kind":           e.Kind.String(),
			"parameter":      e.Parameter,
			"stacktrace":     fmt.Sprintf("%+v", ops),
		})
	} else {
		return log.With(map[string]interface{}{
			"httpStatusCode": httpStatusCode,
			"kind":           e.Kind.String(),
			"parameter":      e.Parameter,
		})
	}
}

func unknownHTTPErrorResponse(ctx *gin.Context, log logger.Logger, e error) {
	er := ErrorResponse{
		Error: ServiceError{
			Message: "Unexpected error - contact support",
		},
	}
	log.WithContext(ctx).Error("Unexpected error: %+v", e)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, er)
}

func newErrResponse(e *Error) ErrorResponse {
	const msg string = "internal server error - please contact support"

	switch e.Kind {
	case Internal, Database:
		return ErrorResponse{
			Error: ServiceError{
				Kind:    Internal.String(),
				Message: msg,
			},
		}
	default:
		return ErrorResponse{
			Error: ServiceError{
				Kind:    e.Kind.String(),
				Param:   string(e.Parameter),
				Message: e.Error(),
			},
		}
	}
}

func httpStatusCode(kind Kind) int {
	switch kind {
	case IO, Internal, Database:
		return http.StatusInternalServerError
	case AlreadyExist:
		return http.StatusConflict
	case NotFound:
		return http.StatusNotFound
	case Validation, BadRequest:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Unauthenticated:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
