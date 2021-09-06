package measurements

import (
    "log"
    "testing"

    "github.com/stretchr/testify/assert"
)

const token = "03e2af2b-0f70-48a9-9a7e-1089781f6e89"

func TestRipe_GetMeasurementResult(t *testing.T) {
    t.Skip(true)

    r := &Ripe{}
    err := r.NewClient(token)
    assert.Nil(t, err)

    // creatred online 32148976
    // created with api - 32221571, 32221572
    // [32221621 32221622]
    got, err := r.GetMeasurementResult(32221571)
    assert.Nil(t, err)
    assert.NotNil(t,got)
    assert.GreaterOrEqual(t, len(got), 5)
    assert.NotNil(t, got[0].Measurement)
    assert.NotNil(t, got[0].Probe)
}

func TestRipe_CreatePing(t *testing.T) {
    t.Skip(true)

    r := &Ripe{}
    err := r.NewClient(token)
    assert.Nil(t, err)

     miners  := []Miner{
         {Address: "x1234", Ip: []string{
                         "213.13.146.142",
                        "143.204.98.83",
         }},
     }


    got, err := r.CreatePing(miners)
    assert.Nil(t, err)
    assert.NotNil(t,got)
    assert.NotNil(t,got.Measurements)
    log.Println(got.Measurements)
}