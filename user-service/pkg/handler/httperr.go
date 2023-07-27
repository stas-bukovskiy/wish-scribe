package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stas-bukovskiy/wish-scribe/user-service/pkg/errs"
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

func newHTTPErrorResponse(ctx *gin.Context, err error) {
	if err == nil {
		nilHTTPErrorResponse(ctx)
	}

	var e *errs.Error
	if errors.As(err, &e) {
		switch e.Kind {
		case errs.Unanticipated:
			unauthenticatedHTTPErrorResponse(ctx, e)
			return
		case errs.Unauthorized:
			unauthorizedHTTPErrorResponse(ctx, e)
			return
		default:
			typicalHTTPErrorResponse(ctx, e)
			return
		}
	}

	unknownHTTPErrorResponse(ctx, err)
}

func nilHTTPErrorResponse(ctx *gin.Context) {
	logrus.Error("nil error - no response body sent")
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

func unauthenticatedHTTPErrorResponse(ctx *gin.Context, e *errs.Error) {
	logErrorResponse("Unauthenticated Request", http.StatusUnauthorized, e)
	ctx.AbortWithStatus(http.StatusUnauthorized)
}

func unauthorizedHTTPErrorResponse(ctx *gin.Context, e *errs.Error) {
	logErrorResponse("Unauthorized Request", http.StatusForbidden, e)
	ctx.AbortWithStatus(http.StatusForbidden)
}

func typicalHTTPErrorResponse(ctx *gin.Context, e *errs.Error) {
	const op errs.Operation = "handler/typicalHTTPErrorResponse"
	const errMsg = "error response sent to client"
	httpStatusCode := httpStatusCode(e.Kind)

	if e.IsEmpty() {
		logrus.Errorf("error sent to %s but with empty body", op)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	logErrorResponse(errMsg, httpStatusCode, e)

	errorResponse := newErrResponse(e)
	ctx.AbortWithStatusJSON(httpStatusCode, errorResponse)
}

func logErrorResponse(msg string, httpStatusCode int, e *errs.Error) {
	if ops := errs.OpStack(e); len(ops) > 0 {
		logrus.WithFields(logrus.Fields{
			"httpStatusCode": httpStatusCode,
			"kind":           e.Kind.String(),
			"parameter":      e.Parameter,
			"stacktrace":     fmt.Sprintf("%+v", ops),
		}).Info(msg)
	} else {
		logrus.WithFields(logrus.Fields{
			"httpStatusCode": httpStatusCode,
			"kind":           e.Kind.String(),
			"parameter":      e.Parameter,
		}).Info(msg)
	}
}

func unknownHTTPErrorResponse(ctx *gin.Context, e error) {
	er := ErrorResponse{
		Error: ServiceError{
			Kind:    errs.Unanticipated.String(),
			Message: "Unexpected error - contact support",
		},
	}
	logrus.Errorf("Unxpected error: %+v", e)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, er)
}

func newErrResponse(e *errs.Error) ErrorResponse {
	const msg string = "internal server error - please contact support"

	switch e.Kind {
	case errs.Internal, errs.Database:
		return ErrorResponse{
			Error: ServiceError{
				Kind:    errs.Internal.String(),
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

func httpStatusCode(kind errs.Kind) int {
	switch kind {
	case errs.Invalid, errs.Exist, errs.NotExist, errs.Private, errs.BrokenLink, errs.Validation, errs.InvalidRequest:
		return http.StatusBadRequest
	case errs.Other, errs.IO, errs.Internal, errs.Database, errs.Unanticipated:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
