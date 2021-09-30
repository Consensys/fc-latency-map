package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeoLocation_updateTrigonometryLatLong(t *testing.T) {
	type fields struct {
		Latitude     float64
		Longitude    float64
		CosLatitude  float64
		SinLatitude  float64
		CosLongitude float64
		SinLongitude float64
	}
	tests := []struct {
		name   string
		fields fields
		want   *GeoLocation
	}{
		{name: "update cos/sin 0,0", fields: fields{
			Latitude:  0,
			Longitude: 0,
		}, want: &GeoLocation{
			Latitude:    0,
			CosLatitude: 1,
			SinLatitude: 0,

			Longitude:    0,
			CosLongitude: 1,
			SinLongitude: 0,
		}},
		{name: "update cos/sin 90,0", fields: fields{
			Latitude:  90,
			Longitude: 0,
		}, want: &GeoLocation{
			Latitude:    90,
			CosLatitude: 6.123233995736757e-17, // almost 0,
			SinLatitude: 1,

			Longitude:    0,
			CosLongitude: 1,
			SinLongitude: 0,
		}},
		{name: "update cos/sin 0,90", fields: fields{
			Latitude:  0,
			Longitude: 90,
		}, want: &GeoLocation{
			Latitude:    0,
			CosLatitude: 1,
			SinLatitude: 0,

			Longitude:    90,
			CosLongitude: 6.123233995736757e-17, // almost 0,
			SinLongitude: 1,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GeoLocation{
				Latitude:  tt.fields.Latitude,
				Longitude: tt.fields.Longitude,
			}
			m.updateTrigonometryLatLong()
			assert.Equal(t, tt.want, m)
		})
	}
}
