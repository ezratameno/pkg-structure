package main

import (
	"reflect"
	"testing"
)

func Test_getImports(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test single import",
			args: args{
				lines: []string{
					"package main",
					`import "fmt"`,
					`func main() {`,
					`    fmt.Println("hello world")`,
					`}`,
				},
			},
			want: []string{"fmt"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getImports(tt.args.lines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getImports() = %v, want %v", got, tt.want)
			}
		})
	}
}
