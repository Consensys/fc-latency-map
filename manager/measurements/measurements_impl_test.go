package measurements

import (
    "reflect"
    "testing"

    "github.com/golang/mock/gomock"
    atlas "github.com/keltia/ripe-atlas"
    "github.com/spf13/viper"
    "github.com/stretchr/testify/assert"

    "github.com/ConsenSys/fc-latency-map/manager/config"
    "github.com/ConsenSys/fc-latency-map/manager/db"
    "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
    "github.com/ConsenSys/fc-latency-map/manager/models"
)

func TestMeasurementServiceImpl_CreateMeasurement(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    type fields struct {
        Conf  *viper.Viper
        DBMgr db.DatabaseMgr
        FMgr  filecoinmgr.FilecoinMgr
    }
    type args struct {
        mr []*atlas.Measurement
    }
    tests := []struct {
        name   string
        fields fields
        args   args
        want   int
    }{
        {name: "empty MeasurementResult", fields: struct {
            Conf  *viper.Viper
            DBMgr db.DatabaseMgr
            FMgr  filecoinmgr.FilecoinMgr
        }{
            Conf:  config.NewMockConfig(),
            DBMgr: db.NewMockDatabaseMgr(),
            FMgr:  filecoinmgr.NewMockFilecoinMgr(ctrl),
        },
            args: struct{ mr []*atlas.Measurement }{
                mr: []*atlas.Measurement{},
            },
            want: 0,
        },
        {name: "not empty MeasurementResult", fields: struct {
            Conf  *viper.Viper
            DBMgr db.DatabaseMgr
            FMgr  filecoinmgr.FilecoinMgr
        }{
            Conf:  config.NewMockConfig(),
            DBMgr: db.NewMockDatabaseMgr(),
            FMgr:  filecoinmgr.NewMockFilecoinMgr(ctrl),
        },
            args: struct{ mr []*atlas.Measurement }{
                mr: []*atlas.Measurement{{
                    Af: 4,
                    ID: 111,
                }},
            },
            want: 1,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := &measurementServiceImpl{
                Conf:  tt.fields.Conf,
                DBMgr: tt.fields.DBMgr,
                FMgr:  tt.fields.FMgr,
            }
            m.UpsertMeasurements(tt.args.mr)
            var rows []*models.Measurement
            err := tt.fields.DBMgr.GetDB().Find(&rows).Error
            assert.Nil(t, err)
            assert.Equal(t, tt.want, len(rows))
        })
    }
}

func TestMeasurementServiceImpl_ImportMeasurement(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    type fields struct {
        Conf  *viper.Viper
        DBMgr db.DatabaseMgr
        FMgr  filecoinmgr.FilecoinMgr
    }
    type args struct {
        mr []atlas.MeasurementResult
    }
    tests := []struct {
        name   string
        fields fields
        args   args
        want   int
    }{
        {name: "empty MeasurementResult", fields: struct {
            Conf  *viper.Viper
            DBMgr db.DatabaseMgr
            FMgr  filecoinmgr.FilecoinMgr
        }{
            Conf:  config.NewMockConfig(),
            DBMgr: db.NewMockDatabaseMgr(),
            FMgr:  filecoinmgr.NewMockFilecoinMgr(ctrl),
        },
            args: struct{ mr []atlas.MeasurementResult }{
                mr: []atlas.MeasurementResult{},
            },
            want: 0,
        },
        {name: "not empty MeasurementResult", fields: struct {
            Conf  *viper.Viper
            DBMgr db.DatabaseMgr
            FMgr  filecoinmgr.FilecoinMgr
        }{
            Conf:  config.NewMockConfig(),
            DBMgr: db.NewMockDatabaseMgr(),
            FMgr:  filecoinmgr.NewMockFilecoinMgr(ctrl),
        },
            args: struct{ mr []atlas.MeasurementResult }{
                mr: []atlas.MeasurementResult{
                    {
                        DstAddr:   "1.2.3.4",
                        MsmID:     123,
                        PrbID:     321,
                        Timestamp: 6546546,
                        Avg:       5,
                        Max:       9,
                        Min:       1,
                    },
                    {
                        DstAddr:   "1.2.3.4",
                        MsmID:     1231,
                        PrbID:     321,
                        Timestamp: 6546546,
                        Avg:       5,
                        Max:       9,
                        Min:       1,
                    },
                },
            },
            want: 2,
        },
        {name: "duplicated measurement id", fields: struct {
            Conf  *viper.Viper
            DBMgr db.DatabaseMgr
            FMgr  filecoinmgr.FilecoinMgr
        }{
            Conf:  config.NewMockConfig(),
            DBMgr: db.NewMockDatabaseMgr(),
            FMgr:  filecoinmgr.NewMockFilecoinMgr(ctrl),
        },
            args: struct{ mr []atlas.MeasurementResult }{
                mr: []atlas.MeasurementResult{
                    {
                        DstAddr:   "1.2.3.4",
                        MsmID:     123,
                        PrbID:     321,
                        Timestamp: 6546546,
                        Avg:       5,
                        Max:       9,
                        Min:       1,
                    },
                    {
                        DstAddr:   "1.2.3.4",
                        MsmID:     123,
                        PrbID:     321,
                        Timestamp: 6546546,
                        Avg:       5,
                        Max:       9,
                        Min:       1,
                    },
                },
            },
            want: 1,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := &measurementServiceImpl{
                Conf:  tt.fields.Conf,
                DBMgr: tt.fields.DBMgr,
                FMgr:  tt.fields.FMgr,
            }
            m.ImportMeasurement(tt.args.mr)
            var rows []*models.MeasurementResult
            err := tt.fields.DBMgr.GetDB().Find(&rows).Error
            assert.Nil(t, err)
            assert.Equal(t, tt.want, len(rows))
        })
    }
}

func TestMeasurementServiceImpl_GetMiners(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    type fields struct {
        Conf  *viper.Viper
        DBMgr db.DatabaseMgr
        FMgr  filecoinmgr.FilecoinMgr
    }
    tests := []struct {
        name   string
        fields fields
        want   []*models.Miner
        create []models.Miner
    }{
        {name: "not empty ", fields: struct {
            Conf  *viper.Viper
            DBMgr db.DatabaseMgr
            FMgr  filecoinmgr.FilecoinMgr
        }{
            Conf:  config.NewMockConfig(),
            DBMgr: db.NewMockDatabaseMgr(),
            FMgr:  filecoinmgr.NewMockFilecoinMgr(ctrl),
        },
            create: []models.Miner{{
                Address:   "f1111",
                IP:        "1.1.1.1",
                Latitude:  1,
                Longitude: 11,
            }},
            want: []*models.Miner{{
                Address:   "f1111",
                IP:        "1.1.1.1",
                Latitude:  1,
                Longitude: 11,
            }},
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := newMeasurementServiceImpl(
                tt.fields.Conf,
                tt.fields.DBMgr,
                tt.fields.FMgr,
            )
            tt.fields.DBMgr.GetDB().Create(tt.create)
            got := m.GetMinersWithGeolocation()
            if len(got) != len(tt.want) {
                t.Errorf("GetMinersWithGeolocation() = %v, want %v", len(got), len(tt.want))
            }
            if len(tt.want) == 0 {
                return
            }
            for i, v := range got {
                assert.Equal(t, tt.want[i].Address, v.Address)
                assert.Equal(t, tt.want[i].IP, v.IP)
                assert.Equal(t, tt.want[i].Latitude, v.Latitude)
                assert.Equal(t, tt.want[i].Longitude, v.Longitude)
            }
        })
    }
}

func TestMeasurementServiceImpl_GetProbIDs(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    type fields struct {
        Conf  *viper.Viper
        DBMgr db.DatabaseMgr
        FMgr  filecoinmgr.FilecoinMgr
    }
    type args struct {
        lat    float64
        long   float64
        amount int
        places []Place
    }
    tests := []struct {
        name      string
        fields    fields
        args      args
        locations []models.Location
        want      []string
        places    []Place
    }{
        {name: "lat,long = 0,0", fields: struct {
            Conf  *viper.Viper
            DBMgr db.DatabaseMgr
            FMgr  filecoinmgr.FilecoinMgr
        }{
            Conf:  config.NewMockConfig(),
            DBMgr: db.NewMockDatabaseMgr(),
            FMgr:  filecoinmgr.NewMockFilecoinMgr(ctrl),
        },
            args: struct {
                lat    float64
                long   float64
                amount int
                places []Place
            }{lat: 0, long: 0, amount: 2},
            locations: []models.Location{
                {IataCode: "iata",
                    Latitude:  1,
                    Longitude: 11,
                    Probes:    []*models.Probe{{ProbeID: 1, Latitude: 1, Longitude: 11}},
                },
            },

            want: []string{},
        },
        {name: "no found probe", fields: struct {
            Conf  *viper.Viper
            DBMgr db.DatabaseMgr
            FMgr  filecoinmgr.FilecoinMgr
        }{
            Conf:  config.NewMockConfig(),
            DBMgr: db.NewMockDatabaseMgr(),
            FMgr:  filecoinmgr.NewMockFilecoinMgr(ctrl),
        },
            args: struct {
                lat    float64
                long   float64
                amount int
                places []Place
            }{
                places: []Place{{
                    ID:        1,
                    Latitude:  1,
                    Longitude: 1,
                }},
                lat:    1,
                long:   1,
                amount: 2,
            },
            locations: []models.Location{
                {IataCode: "iata",
                    Latitude:  1,
                    Longitude: 11,
                },
            },

            want: []string{},
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.fields.Conf.SetDefault("NEAREST_AIRPORTS", tt.args.amount)
            tt.fields.DBMgr.GetDB().
                Create(tt.locations)

            m := newMeasurementServiceImpl(
                tt.fields.Conf,
                tt.fields.DBMgr,
                tt.fields.FMgr,
            )
            if got := m.getProbIDs(tt.args.places, tt.args.lat, tt.args.long); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("getProbIDs() = %v, want %v", got, tt.want)
            }
        })
    }
}
