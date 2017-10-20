package bench

import (
    "fmt"
    "github.com/billyninja/pgtools/connector"
    "time"
)

type Agg struct {
    Count           uint
    TotalLatency    time.Duration
    MaxLatency      time.Duration
    MinLatency      time.Duration
    AvgLatency      time.Duration
}

type SimReport struct {
    Status            SimStatus
    StartedAt         time.Time
    Eta               time.Time
    FinishedAt        *time.Time
    SimulationParams  *SimParams
    UsedConnector     *connector.Connector

    readCount         uint
    writeCount        uint

    InsertSamples     []*Sample
    FlushSample       []*Sample
    ReadSamples       []*Sample
}

type Sample struct {
    WriteCount      uint
    ReadCount       uint
    Latency         time.Duration
}


func (r *SimReport) Finish() {
    now := time.Now()
    r.FinishedAt = &now
    r.Status = SimFinished
}

func AggAnalysis(samples []*Sample) *Agg {
    var sumLat, maxLat time.Duration
    minLat := time.Second * 10

    for _, sample := range samples {

        sumLat += sample.Latency
        if sample.Latency < minLat {
            minLat = sample.Latency
        }
        if sample.Latency > maxLat {
            maxLat = sample.Latency
        }
    }

    l := uint(len(samples))
    if l == 0 {
        l = 1
    }

    return &Agg{
        Count: l,
        TotalLatency: sumLat,
        AvgLatency: sumLat/time.Duration(l),
        MaxLatency: maxLat,
        MinLatency: minLat,
    }
}

func AnalyticPrint(samples []*Sample, agg *Agg) string {
    out := ""
    for _, smp := range samples {
        if smp.Latency >= agg.AvgLatency {

            if smp.Latency >= agg.MaxLatency {
               out += "\x1b[31;1m x" + smp.Latency.String()
            } else {
               out += "\x1b[31;1m -"
            }

            //out += "\x1b[30;1m-"
            continue
        }

        if smp.Latency < agg.AvgLatency {

            if smp.Latency <= agg.MinLatency {
               out += "\x1b[32;1m o"+ smp.Latency.String()
            } else {
               out += "\x1b[32;1m +"
            }

            //out +=  "\x1b[0m-"
            continue
        }
    }

    out +=  "\x1b[0m"
    return out
}

func (r SimReport) String() string {
    out := fmt.Sprintf(
        "\n\nSimulation Report - Status: %s\nStartedAt: %s",
        r.Status,
        r.StartedAt.Format(REPORT_FMT),
    )
    since := time.Since(r.StartedAt)
    if r.Status == SimRunning {
        out += fmt.Sprintf(" (%s ago)", since)
    }
    if r.Status == SimFinished {
        out += fmt.Sprintf("\nFinished at: %s (%s)",
            r.FinishedAt.Format(REPORT_FMT), since)
    }

    out += fmt.Sprintf("\n\n== Used simulation params: \n\t%s", r.SimulationParams)
    out += fmt.Sprintf("\n\nSample analysis:\n\n")

    insert_analysis := AggAnalysis(r.InsertSamples)
    read_analysis := AggAnalysis(r.ReadSamples)

    out += fmt.Sprintf("\n\nSample analysis:\n\nWrite:\n")
    out += AnalyticPrint(r.InsertSamples, insert_analysis)
    out += "\n\nRead:\n"
    out += AnalyticPrint(r.ReadSamples, read_analysis)
    out += "\n\n"

    return out
}
