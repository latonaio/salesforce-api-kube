package str

import "testing"

func TestToFirstUppercase(t *testing.T) {
	tests := []struct {
		args string
		want string
	}{
		{"get", "Get"},
		{"post", "Post"},
		{"put", "Put"},
		{"delete", "Delete"},
	}
	for i, tt := range tests {
		if got := ToFirstUppercase(tt.args); got != tt.want {
			t.Errorf("#%dToFirstUppercase() = %v, want %v", i, got, tt.want)
		}
	}
}
