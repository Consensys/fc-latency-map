package ripemgr

import (
    "github.com/ConsenSys/fc-latency-map/manager/config"
    "github.com/ConsenSys/fc-latency-map/manager/models"
    "github.com/golang/mock/gomock"

    "github.com/stretchr/testify/assert"
    "gopkg.in/h2non/gock.v1"
    "testing"
)

const ourVersion = "0.5.0"
const apiEndpoint = "https://atlas.ripe.net/api/v2"

func Test_NewRipeImpl(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockConfig := config.NewMockConfig()

    r, err := NewRipeImpl(mockConfig)

    assert.Nil(t, err)
    assert.NotNil(t, r)

}

func TestRipeMgrImpl_GetMeasurement_MissingAPIKey(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockConfig := config.NewMockConfig()

    r, _ := NewRipeImpl(mockConfig)
    probes, err := r.GetMeasurement(1)

    assert.NotNil(t, err)
    assert.Nil(t, probes)

}

func TestRipeMgrImpl_GetProbes_MissingAPIKey(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockConfig := config.NewMockConfig()

    r, _ := NewRipeImpl(mockConfig)
    probes, err := r.GetProbes(map[string]string{})

    assert.NotNil(t, err)
    assert.Nil(t, probes)

}

func TestRipeMgrImpl_CreateMeasurementsEmptyMiner(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockConfig := config.NewMockConfig()
    r, _ := NewRipeImpl(mockConfig)
    probes, err := r.CreateMeasurements([]*models.Miner{}, "1", 1)

    assert.Nil(t, err)
    assert.Nil(t, probes)
}

func TestRipeMgrImpl_CreateMeasurementsMissingProbes(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockConfig := config.NewMockConfig()

    r, _ := NewRipeImpl(mockConfig)
    m, err := r.CreateMeasurements([]*models.Miner{
        {
            Address:   "",
            IP:        "",
            Latitude:  0,
            Longitude: 0,
            Port:      0,
        },
    }, "", 1)

    assert.Nil(t, err)
    assert.Nil(t, m)
}

func TestRipeMgrImpl_CreateMeasurements1WithoutNProbesConfigured(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockConfig := config.NewMockConfig()
    mockConfig.Set("RIPE_REQUESTED_PROBES", "")

    r, _ := NewRipeImpl(mockConfig)
    m, err := r.CreateMeasurements([]*models.Miner{
        {
            Address:   "",
            IP:        "",
            Latitude:  0,
            Longitude: 0,
            Port:      0,
        },
    }, "1", 1)

    assert.NotNil(t, err)
    assert.Nil(t, m)
}

func TestRipeMgrImpl_CreateMeasurementsMissingAPIKey(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockConfig := config.NewMockConfig()
    mockConfig.Set("RIPE_ONE_OFF", "false")

    r, _ := NewRipeImpl(mockConfig)
    m, err := r.CreateMeasurements([]*models.Miner{
        {
            Address:   "",
            IP:        "1.1.2.3",
            Latitude:  1,
            Longitude: 1,
            Port:      0,
        },
    }, "1", 1)

    assert.NotNil(t, err)
    assert.Nil(t, m)
}

func TestRipeMgrImpl_CreateMeasurements(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockConfig := config.NewMockConfig()
    mockConfig.Set("RIPE_ONE_OFF", "false")

    //myurl, _ := url.Parse(apiEndpoint)

    //buf := bytes.NewReader(jr)
    defer gock.Off()
    gock.New("https://atlas.ripe.net").
        Post("/api/v2/measurements/traceroute").
        //MatchParam("key", "changeme").
        //MatchHeaders(map[string]string{
        //    "host":       myurl.Host,
        //    "user-agent": fmt.Sprintf("ripe-atlas/%s", ourVersion),
        //}).
        //  Body(buf).
        Reply(200).
        JSON(map[string]string{"foo": "bar"})

    rp, _ := NewRipeImpl(mockConfig)
    m, err := rp.CreateMeasurements([]*models.Miner{
        {
            Address:   "",
            IP:        "1.1.2.3",
            Latitude:  1,
            Longitude: 1,
            Port:      0,
        },
    }, "1", 1)

    assert.NotNil(t, err)
    assert.Nil(t, m)
}
