package apperrors

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/xerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func Test_appError_LevelInfo(t *testing.T) {
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
		want   *appError
	}{
		{
			name: "set level",
			fields: fields{
				next:       nil,
				logMessage: "error",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
			},
			want: &appError{
				next:       nil,
				logMessage: "error",
				level:      levelInfo,
				code:       0,
				message:    "",
				details:    nil,
			},
		},
		{
			name: "change level",
			fields: fields{
				next:       nil,
				logMessage: "error",
				level:      levelError,
				code:       0,
				message:    "",
				details:    nil,
			},
			want: &appError{
				next:       nil,
				logMessage: "error",
				level:      levelInfo,
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
			e := &appError{
				next:       tt.fields.next,
				logMessage: tt.fields.logMessage,
				level:      tt.fields.level,
				frame:      tt.fields.frame,
				code:       tt.fields.code,
				message:    tt.fields.message,
				details:    tt.fields.details,
			}

			got := e.LevelInfo()
			if diff := cmp.Diff(tt.want, got, ignoreFieldsOpt, allowUnexportedOpt); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_appError_LevelError(t *testing.T) {
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
		want   *appError
	}{
		{
			name: "set level",
			fields: fields{
				next:       nil,
				logMessage: "error",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
			},
			want: &appError{
				next:       nil,
				logMessage: "error",
				level:      levelError,
				code:       0,
				message:    "",
				details:    nil,
			},
		},
		{
			name: "change level",
			fields: fields{
				next:       nil,
				logMessage: "error",
				level:      levelInfo,
				code:       0,
				message:    "",
				details:    nil,
			},
			want: &appError{
				next:       nil,
				logMessage: "error",
				level:      levelError,
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
			e := &appError{
				next:       tt.fields.next,
				logMessage: tt.fields.logMessage,
				level:      tt.fields.level,
				frame:      tt.fields.frame,
				code:       tt.fields.code,
				message:    tt.fields.message,
				details:    tt.fields.details,
			}

			got := e.LevelError()
			if diff := cmp.Diff(tt.want, got, ignoreFieldsOpt, allowUnexportedOpt); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_appError_LevelCritical(t *testing.T) {
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
		want   *appError
	}{
		{
			name: "set level",
			fields: fields{
				next:       nil,
				logMessage: "error",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
			},
			want: &appError{
				next:       nil,
				logMessage: "error",
				level:      levelCritical,
				code:       0,
				message:    "",
				details:    nil,
			},
		},
		{
			name: "change level",
			fields: fields{
				next:       nil,
				logMessage: "error",
				level:      levelError,
				code:       0,
				message:    "",
				details:    nil,
			},
			want: &appError{
				next:       nil,
				logMessage: "error",
				level:      levelCritical,
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
			e := &appError{
				next:       tt.fields.next,
				logMessage: tt.fields.logMessage,
				level:      tt.fields.level,
				frame:      tt.fields.frame,
				code:       tt.fields.code,
				message:    tt.fields.message,
				details:    tt.fields.details,
			}

			got := e.LevelCritical()
			if diff := cmp.Diff(tt.want, got, ignoreFieldsOpt, allowUnexportedOpt); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_appError_IsLevelInfo(t *testing.T) {
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
		want   bool
	}{
		{
			name: "simple",
			fields: fields{
				logMessage: "error",
				level:      levelInfo,
				code:       0,
				message:    "",
				details:    nil,
				next:       nil,
			},
			want: true,
		},
		{
			name: "multi layer",
			fields: fields{
				logMessage: "error",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next: &appError{
					logMessage: "error",
					level:      levelInfo,
					code:       0,
					message:    "",
					details:    nil,
					next:       nil,
				},
			},
			want: true,
		},
		{
			name: "missing",
			fields: fields{
				logMessage: "error",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next: &appError{
					logMessage: "error",
					level:      levelError,
					code:       0,
					message:    "",
					details:    nil,
					next:       nil,
				},
			},
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

			got := e.IsLevelInfo()
			if diff := cmp.Diff(tt.want, got); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_appError_IsLevelError(t *testing.T) {
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
		want   bool
	}{
		{
			name: "simple",
			fields: fields{
				logMessage: "error",
				level:      levelError,
				code:       0,
				message:    "",
				details:    nil,
				next:       nil,
			},
			want: true,
		},
		{
			name: "multi layer",
			fields: fields{
				logMessage: "error",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next: &appError{
					logMessage: "error",
					level:      levelError,
					code:       0,
					message:    "",
					details:    nil,
					next:       nil,
				},
			},
			want: true,
		},
		{
			name: "missing",
			fields: fields{
				logMessage: "error",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next: &appError{
					logMessage: "error",
					level:      levelInfo,
					code:       0,
					message:    "",
					details:    nil,
					next:       nil,
				},
			},
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

			got := e.IsLevelError()
			if diff := cmp.Diff(tt.want, got); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_appError_IsLevelCritical(t *testing.T) {
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
		want   bool
	}{
		{
			name: "simple",
			fields: fields{
				logMessage: "error",
				level:      levelCritical,
				code:       0,
				message:    "",
				details:    nil,
				next:       nil,
			},
			want: true,
		},
		{
			name: "multi layer",
			fields: fields{
				logMessage: "error",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next: &appError{
					logMessage: "error",
					level:      levelCritical,
					code:       0,
					message:    "",
					details:    nil,
					next:       nil,
				},
			},
			want: true,
		},
		{
			name: "missing",
			fields: fields{
				logMessage: "error",
				level:      "",
				code:       0,
				message:    "",
				details:    nil,
				next: &appError{
					logMessage: "error",
					level:      levelError,
					code:       0,
					message:    "",
					details:    nil,
					next:       nil,
				},
			},
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

			got := e.IsLevelCritical()
			if diff := cmp.Diff(tt.want, got); len(diff) != 0 {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}
