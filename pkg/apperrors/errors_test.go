package apperrors

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/xerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func Test_appError_Error(t *testing.T) {
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
		name   string
		fields fields
		want   string
	}{
		{
			name: "logMessage exists",
			fields: fields{
				logMessage: "error1",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next: &appError{
					logMessage: "error2",
					level:      levelInfo,
					code:       codes.InvalidArgument,
					message:    "",
					details:    nil,
					next: &appError{
						next:       nil,
						logMessage: "error1",
						level:      "",
						code:       0,
						message:    "",
						details:    nil,
					},
				},
			},
			want: "error1",
		},
		{
			name: "logMessage is empty",
			fields: fields{
				logMessage: "error1",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next: &appError{
					logMessage: "",
					level:      levelInfo,
					code:       codes.InvalidArgument,
					message:    "",
					details:    nil,
					next:       nil,
				},
			},
			want: "no message",
		},
		{
			name: "no appError",
			fields: fields{
				logMessage: "error1",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next:       errors.New("error2"),
			},
			want: "error2",
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

			got := e.Error()
			if diff := cmp.Diff(tt.want, got); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_appError_Is(t *testing.T) {
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
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "InvalidArgument",
			fields: fields{
				logMessage: "error",
				level:      levelInfo,
				code:       codes.InvalidArgument,
				message:    "InvalidArgument",
				details:    nil,
				next:       nil,
			},
			args: args{err: InvalidParameter},
			want: true,
		},
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
					message:    "Internal",
					details:    nil,
					next:       nil,
				},
			},
			args: args{err: InternalServerError},
			want: true,
		},
		{
			name: "not Internal",
			fields: fields{
				logMessage: "error3",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next: &appError{
					logMessage: "error2",
					level:      levelInfo,
					code:       codes.Canceled,
					message:    "Canceled",
					details:    nil,
					next: &appError{
						next:       nil,
						logMessage: "error1",
						level:      levelError,
						code:       codes.NotFound,
						message:    "NotFound",
						details:    nil,
					},
				},
			},
			args: args{err: InternalServerError},
			want: false,
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

			got := e.Is(tt.args.err)
			if diff := cmp.Diff(tt.want, got); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		msg string
	}

	tests := []struct {
		name string
		args args
		want *appError
	}{
		{
			name: "normal",
			args: args{msg: "error"},
			want: &appError{
				next:       nil,
				logMessage: "error",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
			},
		},
	}

	ignoreFieldsOpt := cmpopts.IgnoreFields(appError{}, "frame")
	allowUnexportedOpt := cmp.AllowUnexported(appError{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.msg)
			if diff := cmp.Diff(tt.want, got, ignoreFieldsOpt, allowUnexportedOpt); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}

	tests := []struct {
		name string
		args args
		want *appError
	}{
		{
			name: "format and args",
			args: args{format: "error %d %s", args: []interface{}{1, "str"}},
			want: &appError{
				next:       nil,
				logMessage: "error 1 str",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
			},
		},
		{
			name: "no args",
			args: args{format: "error"},
			want: &appError{
				next:       nil,
				logMessage: "error",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
			},
		},
	}

	ignoreFieldsOpt := cmpopts.IgnoreFields(appError{}, "frame")
	allowUnexportedOpt := cmp.AllowUnexported(appError{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Errorf(tt.args.format, tt.args.args...)
			if diff := cmp.Diff(tt.want, got, ignoreFieldsOpt, allowUnexportedOpt); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	type args struct {
		err error
		msg []string
	}

	tests := []struct {
		name string
		args args
		want *appError
	}{
		{
			name: "wrap normal error",
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
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next: &appError{
					next:       nil,
					logMessage: "error1",
					level:      "",
					code:       0,
					message:    "",
					details:    nil,
				},
			},
		},
		{
			name: "wrap InvalidParameter error",
			args: args{
				err: InvalidParameter.Wrap(New("error1"), "error2"),
				msg: []string{"error3"},
			},
			want: &appError{
				logMessage: "error3",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next: &appError{
					logMessage: "error2",
					level:      levelInfo,
					code:       codes.InvalidArgument,
					message:    "",
					details:    nil,
					next: &appError{
						next:       nil,
						logMessage: "error1",
						level:      "",
						code:       0,
						message:    "",
						details:    nil,
					},
				},
			},
		},
	}

	ignoreFieldsOpt := cmpopts.IgnoreFields(appError{}, "frame")
	allowUnexportedOpt := cmp.AllowUnexported(appError{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Wrap(tt.args.err, tt.args.msg...)
			if diff := cmp.Diff(tt.want, got, ignoreFieldsOpt, allowUnexportedOpt); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestAsAppError(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want *appError
	}{
		{
			name: "New",
			args: args{err: New("error")},
			want: &appError{
				next:       nil,
				logMessage: "error",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
			},
		},
		{
			name: "InvalidParameter",
			args: args{err: InvalidParameter},
			want: &appError{
				next:       nil,
				logMessage: "",
				level:      levelInfo,
				code:       codes.InvalidArgument,
				message:    "",
				details:    nil,
			},
		},
		{
			name: "nil",
			args: args{err: nil},
			want: nil,
		},
		{
			name: "not appError",
			args: args{err: errors.New("error")},
			want: nil,
		},
	}

	ignoreFieldsOpt := cmpopts.IgnoreFields(appError{}, "frame")
	allowUnexportedOpt := cmp.AllowUnexported(appError{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AsAppError(tt.args.err)
			if diff := cmp.Diff(tt.want, got, ignoreFieldsOpt, allowUnexportedOpt); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}
