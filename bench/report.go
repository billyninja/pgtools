package bench

import (
    "fmt"
    "github.com/billyninja/pgtools/connector"
    "github.com/billyninja/pgtools/colors"
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
    FlushSamples       []*Sample
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
               out += colors.Red("x")
            } else {
               out += colors.Yellow("-")
            }
            continue
        }

        if smp.Latency < agg.AvgLatency {

            if smp.Latency <= agg.MinLatency {
               out += colors.Green("o")
            } else {
               out += colors.Blue("+")
            }
            continue
        }
    }

    out += "\nCaption:"
    out += "\n" + colors.Green("o") +  " - lowest latency      - " + colors.Green(agg.MinLatency.String())
    out += "\n" + colors.Blue("+") +  " - bellow avg. latency - " + colors.Blue("< " + agg.AvgLatency.String())
    out += "\n" + colors.Yellow("-") + " - above avg. latency  - " + colors.Yellow("> " + agg.AvgLatency.String())
    out += "\n" + colors.Red("x") +    " - highest latency     - " + colors.Red(agg.MaxLatency.String())
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

    ips := float64(r.writeCount)/since.Seconds()
    out += fmt.Sprintf("\nAt: %.2f inserts/s on avarage", ips)
    out += fmt.Sprintf("\n\n== Used simulation params: \n\t%s", r.SimulationParams)

    insert_analysis := AggAnalysis(r.InsertSamples)
    flush_analysis := AggAnalysis(r.FlushSamples)
    read_analysis := AggAnalysis(r.ReadSamples)

    out += fmt.Sprintf("\n\nSample analysis:\n\nWrite:\n")
    out += AnalyticPrint(r.InsertSamples, insert_analysis)
    out += "\n\n" + colors.Bold("Flush") + ":\n"
    out += AnalyticPrint(r.FlushSamples, flush_analysis)
    out += "\n\nRead:\n"
    out += AnalyticPrint(r.ReadSamples, read_analysis)
    out += "\n\n"

    return out
}
