package errs

import (
	"errors"
	"fmt"
	errors2 "github.com/pkg/errors"
	"runtime"
	"sort"
)

type Kind uint8
type Operation string
type Parameter string

const (
	IO           Kind = iota // Unclassified error. This value is not printed in the error message.
	AlreadyExist             // Item already exists.
	NotFound                 // Item does not exist.
	Internal                 // Internal error or inconsistency.
	Database                 // Error from database.
	Validation               // Input validation error.
	BadRequest               // Invalid Request
	Unauthorized
	Unauthenticated
)

type Error struct {
	Operation Operation
	Kind      Kind
	Parameter Parameter
	Err       error
}

func NewError(args ...interface{}) error {
	type stackTracer interface {
		StackTrace() errors2.StackTrace
	}

	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Operation:
			e.Operation = arg
		case Kind:
			e.Kind = arg
		case string:
			e.Err = errors2.New(arg)
		case *Error:
			errorCopy := *arg
			e.Err = &errorCopy
		case error:
			_, ok := arg.(stackTracer)
			if ok {
				e.Err = arg
			} else {
				e.Err = errors2.New(arg.Error())
			}
		case Parameter:
			e.Parameter = arg
		default:
			_, file, line, _ := runtime.Caller(1)
			return fmt.Errorf("errs.NewError: bad call from %s:%d: %v, unknown type %T, value %v in error call", file, line, args, arg, arg)
		}
	}

	prev, ok := e.Err.(*Error)
	if !ok {
		return e
	}

	if prev.Parameter == e.Parameter {
		prev.Parameter = ""
	}
	if e.Parameter == "" {
		e.Parameter = prev.Parameter
		prev.Parameter = ""
	}

	return e
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) IsEmpty() bool {
	return e.Kind == 0 && e.Parameter == "" && e.Err == nil
}

func OpStack(e error) []string {
	type o struct {
		Op    string
		Order int
	}
	var os []o

	currentErr := e
	i := 0
	for errors.Unwrap(currentErr) != nil {
		var errsErr *Error
		if errors.As(e, &errsErr) {
			if errsErr.Operation != "" {
				op := o{string(errsErr.Operation), i}
				os = append(os, op)
			}
		}
		currentErr = errors.Unwrap(currentErr)
		i++
	}

	sort.Slice(os, func(i, j int) bool { return os[i].Order > os[j].Order })

	var ops []string
	for _, op := range os {
		ops = append(ops, op.Op)
	}

	return ops
}

func (k Kind) String() string {
	switch k {
	case IO:
		return "I/O error"
	case AlreadyExist:
		return "item already exists"
	case NotFound:
		return "item does not exist"
	case Internal:
		return "internal error"
	case Database:
		return "database error"
	case Validation:
		return "input validation error"
	case BadRequest:
		return "invalid request error"
	case Unauthorized:
		return "unauthorized request"
	case Unauthenticated:
		return "unauthenticated request"
	}
	return "unknown error kind"
}

func KindIs(kind Kind, err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return kind == e.Kind
	}
	return false
}
