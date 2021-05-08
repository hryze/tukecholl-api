package apperrors

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/xerrors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func Test_appError_Wrap(t *testing.T) {
	type fields struct {
		next       error
		logMessage string
		level      level
		frame      xerrors.Frame
		code       codes.Code
		message    string
		details    []proto.Message
	}

	type args struct {
		err error
		msg []string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *appError
	}{
		{
			name: "exists message in argument",
			fields: fields{
				logMessage: "",
				level:      levelInfo,
				code:       codes.InvalidArgument,
				message:    "invalid argument",
				details:    nil,
				next:       nil,
			},
			args: args{
				err: &appError{
					logMessage: "error1",
					level:      "",
					code:       0,
					message:    "",
					details:    nil,
					next:       nil,
				},
				msg: []string{"error2"},
			},
			want: &appError{
				logMessage: "error2",
				level:      levelInfo,
				code:       codes.InvalidArgument,
				message:    "invalid argument",
				details:    nil,
				next: &appError{
					logMessage: "error1",
					level:      "",
					code:       0,
					message:    "",
					details:    nil,
					next:       nil,
				},
			},
		},
		{
			name: "missing message in argument",
			fields: fields{
				logMessage: "",
				level:      levelInfo,
				code:       codes.InvalidArgument,
				message:    "invalid argument",
				details:    nil,
				next:       nil,
			},
			args: args{
				err: &appError{
					logMessage: "error1",
					level:      "",
					code:       0,
					message:    "",
					details:    nil,
					next:       nil,
				},
				msg: []string{},
			},
			want: &appError{
				logMessage: "InvalidArgument",
				level:      levelInfo,
				code:       codes.InvalidArgument,
				message:    "invalid argument",
				details:    nil,
				next: &appError{
					logMessage: "error1",
					level:      "",
					code:       0,
					message:    "",
					details:    nil,
					next:       nil,
				},
			},
		},
	}

	ignoreFieldsOpt := cmpopts.IgnoreFields(appError{}, "frame")
	allowUnexportedOpt := cmp.AllowUnexported(appError{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := appError{
				next:       tt.fields.next,
				logMessage: tt.fields.logMessage,
				level:      tt.fields.level,
				frame:      tt.fields.frame,
				code:       tt.fields.code,
				message:    tt.fields.message,
				details:    tt.fields.details,
			}

			got := e.Wrap(tt.args.err, tt.args.msg...)
			if diff := cmp.Diff(tt.want, got, ignoreFieldsOpt, allowUnexportedOpt); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_appError_Status(t *testing.T) {
	badRequestDetails := &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			{
				Field:       "user id",
				Description: "invalid user id",
			},
		},
	}

	invalidArgumentStatus, _ := status.New(codes.InvalidArgument, "invalid argument").WithDetails(badRequestDetails)

	type fields struct {
		next       error
		logMessage string
		level      level
		frame      xerrors.Frame
		code       codes.Code
		message    string
		details    []proto.Message
	}

	tests := []struct {
		name    string
		fields  fields
		want    *status.Status
		wantErr bool
	}{
		{
			name: "Internal",
			fields: fields{
				logMessage: "error2",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next: &appError{
					logMessage: "error1",
					level:      levelError,
					code:       codes.Internal,
					message:    "internal",
					details:    nil,
					next:       nil,
				},
			},
			want:    status.New(codes.Internal, "internal"),
			wantErr: false,
		},
		{
			name: "InvalidArgument with details",
			fields: fields{
				logMessage: "error2",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next: &appError{
					logMessage: "error1",
					level:      levelInfo,
					code:       codes.InvalidArgument,
					message:    "invalid argument",
					details: []proto.Message{
						&errdetails.BadRequest{
							FieldViolations: []*errdetails.BadRequest_FieldViolation{
								{
									Field:       "user id",
									Description: "invalid user id",
								},
							},
						},
					},
					next: nil,
				},
			},
			want:    invalidArgumentStatus,
			wantErr: false,
		},
		{
			name: "no status",
			fields: fields{
				logMessage: "error1",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next:       nil,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &appError{
				next:       tt.fields.next,
				logMessage: tt.fields.logMessage,
				level:      tt.fields.level,
				frame:      tt.fields.frame,
				code:       tt.fields.code,
				message:    tt.fields.message,
				details:    tt.fields.details,
			}

			got, err := e.Status()
			if diff := cmp.Diff(tt.wantErr, err != nil); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Status() want = %v, got %v", tt.want, got)
			}
		})
	}
}

func Test_appError_SetMessage(t *testing.T) {
	type fields struct {
		next       error
		logMessage string
		level      level
		frame      xerrors.Frame
		code       codes.Code
		message    string
		details    []proto.Message
	}

	type args struct {
		msg string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *appError
	}{
		{
			name: "normal",
			fields: fields{
				logMessage: "",
				level:      levelError,
				code:       codes.Internal,
				message:    "",
				details:    nil,
				next:       nil,
			},
			args: args{msg: "internal error"},
			want: &appError{
				logMessage: "",
				level:      levelError,
				code:       codes.Internal,
				message:    "internal error",
				details:    nil,
				next:       nil,
			},
		},
	}

	ignoreFieldsOpt := cmpopts.IgnoreFields(appError{}, "frame")
	allowUnexportedOpt := cmp.AllowUnexported(appError{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &appError{
				next:       tt.fields.next,
				logMessage: tt.fields.logMessage,
				level:      tt.fields.level,
				frame:      tt.fields.frame,
				code:       tt.fields.code,
				message:    tt.fields.message,
				details:    tt.fields.details,
			}

			got := e.SetMessage(tt.args.msg)
			if diff := cmp.Diff(tt.want, got, ignoreFieldsOpt, allowUnexportedOpt); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_appError_AddBadRequestFieldViolation(t *testing.T) {
	type fields struct {
		next       error
		logMessage string
		level      level
		frame      xerrors.Frame
		code       codes.Code
		message    string
		details    []proto.Message
	}

	type args struct {
		field string
		desc  string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *appError
	}{
		{
			name: "new",
			fields: fields{
				logMessage: "",
				level:      levelInfo,
				code:       codes.InvalidArgument,
				message:    "",
				details:    nil,
				next:       nil,
			},
			args: args{
				field: "user id",
				desc:  "invalid user id",
			},
			want: &appError{
				logMessage: "",
				level:      levelInfo,
				code:       codes.InvalidArgument,
				message:    "",
				details: []proto.Message{
					&errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequest_FieldViolation{
							{
								Field:       "user id",
								Description: "invalid user id",
							},
						},
					},
				},
				next: nil,
			},
		},
		{
			name: "add",
			fields: fields{
				logMessage: "",
				level:      levelInfo,
				code:       codes.InvalidArgument,
				message:    "",
				details: []proto.Message{
					&errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequest_FieldViolation{
							{
								Field:       "user id",
								Description: "invalid user id",
							},
						},
					},
				},
				next: nil,
			},
			args: args{
				field: "user name",
				desc:  "invalid user name",
			},
			want: &appError{
				logMessage: "",
				level:      levelInfo,
				code:       codes.InvalidArgument,
				message:    "",
				details: []proto.Message{
					&errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequest_FieldViolation{
							{
								Field:       "user id",
								Description: "invalid user id",
							},
							{
								Field:       "user name",
								Description: "invalid user name",
							},
						},
					},
				},
				next: nil,
			},
		},
	}

	ignoreFieldsOpt := cmpopts.IgnoreFields(appError{}, "frame")
	ignoreUnexportedOpt := cmpopts.IgnoreUnexported(errdetails.BadRequest{}, errdetails.BadRequest_FieldViolation{})
	allowUnexportedOpt := cmp.AllowUnexported(appError{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &appError{
				next:       tt.fields.next,
				logMessage: tt.fields.logMessage,
				level:      tt.fields.level,
				frame:      tt.fields.frame,
				code:       tt.fields.code,
				message:    tt.fields.message,
				details:    tt.fields.details,
			}

			got := e.AddBadRequestFieldViolation(tt.args.field, tt.args.desc)
			if diff := cmp.Diff(tt.want, got, ignoreFieldsOpt, ignoreUnexportedOpt, allowUnexportedOpt); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}
