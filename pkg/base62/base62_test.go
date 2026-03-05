package base62

import "testing"

func TestBase62Encode(t *testing.T) {
	type args struct {
		u uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "1", args: args{u: 0}, want: "0"},
		{name: "2", args: args{u: 6349}, want: "1Ep"},
		{name: "3", args: args{u: 62}, want: "10"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base62Encode(tt.args.u); got != tt.want {
				t.Errorf("Base62Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase62Decode1(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{name: "1", args: args{s: "1Ep"}, want: 6349},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base62Decode(tt.args.s); got != tt.want {
				t.Errorf("Base62Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
