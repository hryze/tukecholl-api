package apperrors

import (
	"golang.org/x/xerrors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newInvalidArgument() *appError {
	return &appError{
		level: levelInfo,
		code:  codes.InvalidArgument,
	}
}

func newInternalServerError() *appError {
	return &appError{
		level: levelError,
		code:  codes.Internal,
	}
}

func (e appError) create(msg string) *appError {
	e.logMessage = msg
	e.frame = xerrors.Caller(2)

	return &e
}

func (e appError) Wrap(err error, msg ...string) *appError {
	var m string
	if len(msg) != 0 {
		m = msg[0]
	} else {
		m = e.code.String()
	}

	ne := e.create(m)
	ne.next = err

	return ne
}

func (e *appError) Status() (*status.Status, error) {
	if e.code == codes.OK {
		if next := AsAppError(e.next); next != nil {
			return next.Status()
		}

		return nil, xerrors.New("unknown gRPC code")
	}

	st := status.New(e.code, e.message)
	var err error

	for _, detail := range e.details {
		switch t := detail.(type) {
		case *errdetails.BadRequest:
			if st, err = st.WithDetails(t); err != nil {
				return nil, err
			}
		default:
			return nil, xerrors.New("failed to type assertion from proto.Message")
		}
	}

	return st, nil
}

func (e *appError) SetMessage(msg string) *appError {
	e.message = msg

	return e
}

func (e *appError) AddBadRequestFieldViolation(field, desc string) *appError {
	fieldViolation := &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: desc,
	}

	for i, detail := range e.details {
		badRequestDetail, ok := detail.(*errdetails.BadRequest)
		if !ok {
			continue
		}

		badRequestDetail.FieldViolations = append(badRequestDetail.FieldViolations, fieldViolation)
		e.details[i] = badRequestDetail

		return e
	}

	badRequestDetail := &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			fieldViolation,
		},
	}

	e.details = append(e.details, badRequestDetail)

	return e
}
