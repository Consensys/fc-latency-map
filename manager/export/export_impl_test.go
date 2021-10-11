package export

import (
	"testing"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/spf13/viper"
)

func TestExportServiceImpl_export(t *testing.T) {
	type fields struct {
		Conf  *viper.Viper
		DBMgr db.DatabaseMgr
	}
	type args struct {
		fn string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		createResults   []*models.MeasurementResult
		createLocations []*models.Location
		createMiners    []*models.Miner
	}{
		{name: "empty file name",
			fields: struct {
				Conf  *viper.Viper
				DBMgr db.DatabaseMgr
			}{
				Conf:  config.NewMockConfig(),
				DBMgr: db.NewMockDatabaseMgr(),
			},
			createMiners:  []*models.Miner{{Address: "f11111", IP: "1.1.1.1"}},
			createResults: []*models.MeasurementResult{{MeasurementID: 1, ProbeID: 1, IP: "1.1.1.1"}},
			createLocations: []*models.Location{
				{IataCode: "iata",
					Latitude:  1,
					Longitude: 11,
					Probes: []*models.Probe{{ProbeID: 1, Latitude: 1, Longitude: 11},
						{ProbeID: 3, Latitude: 1, Longitude: 11},
						{ProbeID: 2, Latitude: 1, Longitude: 11},
					},
				},
				{IataCode: "iata1",
					Latitude:  1,
					Longitude: 11,
					Probes: []*models.Probe{{ProbeID: 1, Latitude: 1, Longitude: 11},
						{ProbeID: 4, Latitude: 1, Longitude: 11},
						{ProbeID: 5, Latitude: 1, Longitude: 11},
					},
				},
			},
			args: struct{ fn string }{
				fn: "",
			},
		},
		{name: "miner empty ip",
			fields: struct {
				Conf  *viper.Viper
				DBMgr db.DatabaseMgr
			}{
				Conf:  config.NewMockConfig(),
				DBMgr: db.NewMockDatabaseMgr(),
			},
			createMiners:  []*models.Miner{{Address: "f11111", IP: ""}},
			createResults: []*models.MeasurementResult{{MeasurementID: 1, ProbeID: 1, IP: "1.1.1.1"}},
			createLocations: []*models.Location{
				{
					IataCode:  "iata",
					Latitude:  1,
					Longitude: 11,
					Probes:    []*models.Probe{{ProbeID: 1, Latitude: 1, Longitude: 11}},
				},
			},
			args: struct{ fn string }{
				fn: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbc := tt.fields.DBMgr.GetDB()
			for _, l := range tt.createLocations {
				dbc.Create(l.Probes)
			}
			dbc.Create(tt.createResults)
			dbc.Create(tt.createLocations)
			dbc.Create(tt.createMiners)
			m := &ExportServiceImpl{
				Conf:  tt.fields.Conf,
				DBMgr: tt.fields.DBMgr,
			}
			m.Export(tt.args.fn)
		})
	}
}
