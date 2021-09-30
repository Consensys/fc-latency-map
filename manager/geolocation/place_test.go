package geolocation

import (
	"math"
	"testing"
)

func TestRadians(t *testing.T) {
	type args struct {
		degree float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{name: "test 0 degrees", args: args{0}, want: 0},
		{name: "test 90 degrees", args: args{90}, want: math.Pi / 2},
		{name: "test 180 degrees", args: args{180}, want: math.Pi},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := radians(tt.args.degree); got != tt.want {
				t.Errorf("radians() = %v, want %v", got, tt.want)
			}
		})
	}
}
