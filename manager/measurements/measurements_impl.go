package measurements

import (
    "fmt"
    "log"
    "time"

    "github.com/keltia/ripe-atlas"
)

type Ripe struct {
    c         *atlas.Client
    StartTime int
    StopTime  int
    IsOneOff  bool
}

func (r *Ripe) NewClient(t string, cfgs ...atlas.Config) error {
    if cfgs == nil {
        cfgs = append(cfgs, atlas.Config{
            APIKey: t,
        })
    }

    c, err := atlas.NewClient(cfgs...)
    if err != nil {
        log.Println("Connecting to measurements API", err)
        return err
    }
    r.c = c
    ver := atlas.GetVersion()
    log.Println("api version ", ver)

    return nil
}

func (r *Ripe) GetMeasurement(id int) (m *atlas.Measurement, err error) {
    return r.c.GetMeasurement(id)
}

func (r *Ripe) CreatePing(miners []Miner) (m *atlas.MeasurementResp, err error) {
    var d []atlas.Definition
    for _, miner := range miners {
        for _, ip := range miner.Ip {
            d = append(d, atlas.Definition{
                Description:    fmt.Sprintf("%s ping to %s", miner.Address, ip),
                AF:             4,
                Target:         ip,
                Spread:         4,
                Packets:        2,
                PacketInterval: 10,
                // Tags: []string{
                //     miner.Address, ip,
                // },
                ExtraWait:      5,
                Interval:       6000, // 10 minutes
                Retry:          1,
                UDPPayloadSize: 2,
                Size:           48,
                Type:           "ping",
                GroupID:        0,
                Group:          "",
                InWifiGroup:    false,
            })
        }
    }

    mr := &atlas.MeasurementRequest{
        Definitions: d,
        StartTime:   int(time.Now().Unix()),
        StopTime:    int(time.Now().Unix() + 3600), // 1 hour
        IsOneoff:    false,
        Probes: []atlas.ProbeSet{
            {
                Type:      "area",
                Value:     "WW",
                Requested: 5,
            },
        },
    }

    log.Println("NewMeasurement ", fmt.Sprintf("%#v\n", d))

    return r.c.Ping(mr)
}

func (r *Ripe) GetMeasurementResult(id int) ([]MeasurementResults, error) {
    var probeResults []MeasurementResults
    var probes map[int]*atlas.Probe
    results, err := r.c.GetResults(id)
    if err != nil {
        log.Println("Getting measurements API ", err)
        return nil, err
    }

    // create a str
    for _, ripeResult := range results.Results {
        p, found := probes[ripeResult.PrbID]
        if !found {
            p, err = r.c.GetProbe(ripeResult.PrbID)
            if err != nil {
                log.Println("Getting probe API ", ripeResult.PrbID, err)
                continue
            }
        }

        probeResults = append(probeResults, MeasurementResults{
            Measurement: ripeResult,
            Probe:       *p,
        })
    }
    return probeResults, err
}
