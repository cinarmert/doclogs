package container

import "testing"

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name string
		s    Status
		want string
	}{
		{
			name: "Idle",
			s:    Idle,
			want: "Idle",
		},
		{
			name: "Live",
			s:    Live,
			want: "Live",
		},
		{
			name: "Errored",
			s:    Errored,
			want: "Errored",
		},
		{
			name: "Terminated",
			s:    Terminated,
			want: "Terminated",
		},
		{
			name: "Unknown Status",
			s:    5,
			want: "Unknown Status",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
