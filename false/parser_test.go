package false

import (
	"reflect"
	"sandbox-vm/vm"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "check 1+2 program",
			args: struct {
				str string
			}{
				str: "1 2 \\ + .",
			},
			want:    []int{1, 1, 1, 2, vm.InstrSwap, vm.InstrPlus, vm.InstrWriteInt, vm.InstrEnd},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, err := p.Parse(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
