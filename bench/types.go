package bench

import (
    "fmt"
    "time"
    "github.com/billyninja/pgtools/scanner"
    "github.com/billyninja/pgtools/colors"
)

type SimStatus uint8
type WipeMode uint8
type CountMode uint8
type ReaderFunc uint8
type SimParams struct {
    Table               scanner.TableName
    Wipe                WipeMode
    Count               uint
    CountMode           CountMode
    ReadFunc            ReaderFunc
    InsertsPerSecond    uint
    ReadsPerSecond      uint
    SleepPerInsert      time.Duration
    SleepPerRead        time.Duration
}



const (
    REPORT_FMT              string = "01/02 15:04:05"

    SimNotStarted           SimStatus = iota
    SimRunning
    SimFinished

    WipeNever               WipeMode = iota
    WipeBefore
    WipeAfter
    WipeBeforeAndAfter

    FillIncrement           CountMode = iota
    FillUntil

    ReaderGlobalCount       ReaderFunc = iota
    ReaderUnitarySelect
    ReaderBigSelect
    ReaderBigAgg
)

func (ss SimStatus) String() string {
    switch ss {
        case SimNotStarted:
            return "Not started"
        case SimRunning:
            return "Running"
        case SimFinished:
            return "Finished"
    }

    return "unknown sim. status?!"
}


func (wm WipeMode) String() string {
    switch wm {
        case WipeNever:
            return "Dont wipe the table"
        case WipeBefore:
            return "Wipe the table before running the tests"
        case WipeAfter:
            return "Wipe the table after running the tests"
        case WipeBeforeAndAfter:
            return "Wipe the table before AND after running the tests"
    }

    return "unknown wipe mode?!"
}


func (cm CountMode) String() string {
    switch cm {
        case FillIncrement:
            return "Increment the table by the informed Count parameter"
        case FillUntil:
            return "Increment the table until it reaches the informed Count parameter"
    }

    return "unknown count mode?!"
}


func (rf ReaderFunc) String() string {

    switch rf {
        case ReaderGlobalCount:
            return "Count(*) on the entire table"
        case ReaderUnitarySelect:
            return "Select on a single record"
        case ReaderBigSelect:
            return "Select retrieving a big range of records"
        case ReaderBigAgg:
            return "Select using a Agg function on a big range of records"
    }

    return "unknown wipe mode?!"
}


func (sp SimParams) String() string {
    base := `
    %s:        %s
    %s:            %s
    %s:         %s (%s)
    %s:   %s (%s sleep)
    %s:     %s (%s sleep, %s)
    `
    return fmt.Sprintf(
        base,
        colors.Yellow("SelectedTable"),
        colors.Bold(sp.Table),
        colors.Yellow("Wipe Mode"),
        colors.Bold(sp.Wipe),
        colors.Yellow("Insert Count"),
        colors.Bold(sp.Count),
        sp.CountMode,
        colors.Yellow("Inserts Per Second"),
        colors.Bold(sp.InsertsPerSecond),
        sp.SleepPerInsert,
        colors.Yellow("Reads Per Second"),
        colors.Bold(sp.ReadsPerSecond),
        sp.SleepPerRead,
        sp.ReadFunc)
}

