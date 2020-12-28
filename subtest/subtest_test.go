package subtest

import (
	"testing"
)

func TestCamel(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"camel case", "thisIsACamelCaseString", "this_is_a_camel_case_string"},
		{"with space", "with a space", "with a space"},
		{"ends with A", "endsWithA", "ends_with_a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Camel(tt.arg); got != tt.want {
				t.Fatalf("Camel() = %v, want %v", got, tt.want)
			}
		})
	}
}
