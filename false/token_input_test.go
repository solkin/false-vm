package false

import (
	"false-vm/input"
	"testing"
)

func TestTokenInput_IsInt(t *testing.T) {
	type fields struct {
		Input input.StringInput
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "check 0 is int",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "0"}},
			want:   true,
		},
		{
			name:   "check a is not int",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "a"}},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := &TokenInput{
				Input: tt.fields.Input,
			}
			if got := ti.IsInt(); got != tt.want {
				t.Errorf("IsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenInput_ReadInt(t *testing.T) {
	type fields struct {
		Input input.StringInput
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name:    "check 1 is read as 1",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "1"}},
			want:    1,
			wantErr: false,
		},
		{
			name:    "check 123 is read as 123",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "123"}},
			want:    123,
			wantErr: false,
		},
		{
			name:    "check 12a is read as 12",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "12a"}},
			want:    12,
			wantErr: false,
		},
		{
			name:    "check a12 produces error",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "a12"}},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := &TokenInput{
				Input: tt.fields.Input,
			}
			got, err := ti.ReadInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenInput_IsCharCode(t *testing.T) {
	type fields struct {
		Input input.StringInput
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "check 'a is char",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "'a"}},
			want:   true,
		},
		{
			name:   "check a is not char",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "a"}},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := &TokenInput{
				Input: tt.fields.Input,
			}
			if got := ti.IsCharCode(); got != tt.want {
				t.Errorf("IsCharCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenInput_ReadCharCode(t *testing.T) {
	type fields struct {
		Input input.StringInput
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name:    "check 'a is read as char a",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "'a"}},
			want:    int('a'),
			wantErr: false,
		},
		{
			name:    "check 'aa is read as char a",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "'aa"}},
			want:    int('a'),
			wantErr: false,
		},
		{
			name:    "check aa is not char",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "aa"}},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := &TokenInput{
				Input: tt.fields.Input,
			}
			got, err := ti.ReadCharCode()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadCharCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadCharCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenInput_IsCommand(t *testing.T) {
	type fields struct {
		Input input.StringInput
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "check $ is a stack operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "$"}},
			want:   true,
		},
		{
			name:   "check % is a stack operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "%"}},
			want:   true,
		},
		{
			name:   "check \\ is a stack operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "\\"}},
			want:   true,
		},
		{
			name:   "check @ is a stack operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "@"}},
			want:   true,
		},
		{
			name:   "check ø is a stack operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "ø"}},
			want:   true,
		},
		{
			name:   "check + is an arithmetic operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "+"}},
			want:   true,
		},
		{
			name:   "check - is an arithmetic operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "-"}},
			want:   true,
		},
		{
			name:   "check * is an arithmetic operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "*"}},
			want:   true,
		},
		{
			name:   "check / is an arithmetic operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "/"}},
			want:   true,
		},
		{
			name:   "check _ is an arithmetic operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "_"}},
			want:   true,
		},
		{
			name:   "check & is an arithmetic operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "&"}},
			want:   true,
		},
		{
			name:   "check | is an arithmetic operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "|"}},
			want:   true,
		},
		{
			name:   "check ~ is an arithmetic operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "~"}},
			want:   true,
		},
		{
			name:   "check > is an comparison operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: ">"}},
			want:   true,
		},
		{
			name:   "check = is an comparison operation",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "="}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := &TokenInput{
				Input: tt.fields.Input,
			}
			if got := ti.IsCommand(); got != tt.want {
				t.Errorf("IsCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenInput_IsString(t *testing.T) {
	type fields struct {
		Input input.StringInput
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "check \"a\" is string",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "\"a\""}},
			want:   true,
		},
		{
			name:   "check a is not int",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "a"}},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := &TokenInput{
				Input: tt.fields.Input,
			}
			if got := ti.IsString(); got != tt.want {
				t.Errorf("IsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenInput_ReadString(t *testing.T) {
	type fields struct {
		Input input.StringInput
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "check \"1\" is read as string 1",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "\"1\""}},
			want:    "1",
			wantErr: false,
		},
		{
			name:    "check \"123\" is read as string 123",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "\"123\""}},
			want:    "123",
			wantErr: false,
		},
		{
			name:    "check \"12\". is read as string 12",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "\"12\"a"}},
			want:    "12",
			wantErr: false,
		},
		{
			name:    "check a12\" produces error",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "a12\""}},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := &TokenInput{
				Input: tt.fields.Input,
			}
			got, err := ti.ReadString()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenInput_IsVar(t *testing.T) {
	type fields struct {
		Input input.StringInput
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "check 0 is not var",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "0"}},
			want:   false,
		},
		{
			name:   "check a is var",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "a"}},
			want:   true,
		},
		{
			name:   "check aa is var",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "aa"}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := &TokenInput{
				Input: tt.fields.Input,
			}
			if got := ti.IsVar(); got != tt.want {
				t.Errorf("IsVar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenInput_ReadVar(t *testing.T) {
	type fields struct {
		Input input.StringInput
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		want1   rune
		wantErr bool
	}{
		{
			name:    "check 1 is not var",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "1"}},
			want:    "",
			want1:   0,
			wantErr: true,
		},
		{
			name:    "check 1a is not var",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "1a"}},
			want:    "",
			want1:   0,
			wantErr: true,
		},
		{
			name:    "check a is invalid var mode",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "a"}},
			want:    "",
			want1:   0,
			wantErr: true,
		},
		{
			name:    "check a: is store var a",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "a:"}},
			want:    "a",
			want1:   STORE_VAR,
			wantErr: false,
		},
		{
			name:    "check a;12 is fetch var a",
			fields:  struct{ Input input.StringInput }{Input: input.StringInput{Str: "a;12"}},
			want:    "a",
			want1:   FETCH_VAR,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := &TokenInput{
				Input: tt.fields.Input,
			}
			got, got1, err := ti.ReadVar()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadVar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadVar() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ReadVar() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestTokenInput_IsSubStart(t *testing.T) {
	type fields struct {
		Input input.StringInput
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "check a[ is not sub",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "a["}},
			want:   false,
		},
		{
			name:   "check [] is sub",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "[]"}},
			want:   true,
		},
		{
			name:   "check [[]] is sub",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "[[]]"}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := &TokenInput{
				Input: tt.fields.Input,
			}
			if got := ti.IsSubStart(); got != tt.want {
				t.Errorf("IsSubStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenInput_IsSubEnd(t *testing.T) {
	type fields struct {
		Input input.StringInput
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "check a] is not sub end",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "a]"}},
			want:   false,
		},
		{
			name:   "check ] is sub end",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "]"}},
			want:   true,
		},
		{
			name:   "check ]a is sub end",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "]a"}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := &TokenInput{
				Input: tt.fields.Input,
			}
			if got := ti.IsSubEnd(); got != tt.want {
				t.Errorf("IsSubEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenInput_IsSubCall(t *testing.T) {
	type fields struct {
		Input input.StringInput
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "check a! is not sub call",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "a!"}},
			want:   false,
		},
		{
			name:   "check ! is sub call",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "!"}},
			want:   true,
		},
		{
			name:   "check !a is sub call",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "!a"}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := &TokenInput{
				Input: tt.fields.Input,
			}
			if got := ti.IsSubCall(); got != tt.want {
				t.Errorf("IsSubCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenInput_NextSkipWhitespaces(t *testing.T) {
	type fields struct {
		Input input.StringInput
	}
	tests := []struct {
		name   string
		fields fields
		want   rune
	}{
		{
			name:   "check a returns a",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "a"}},
			want:   'a',
		},
		{
			name:   "check ' aa' returns a",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: " aa"}},
			want:   'a',
		},
		{
			name:   "check 'aa' returns a",
			fields: struct{ Input input.StringInput }{Input: input.StringInput{Str: "aa"}},
			want:   'a',
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := &TokenInput{
				Input: tt.fields.Input,
			}
			if got := ti.NextSkipWhitespaces(); got != tt.want {
				t.Errorf("NextSkipWhitespaces() = %v, want %v", got, tt.want)
			}
		})
	}
}
