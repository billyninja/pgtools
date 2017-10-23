package bench

import (
    "fmt"
    "github.com/billyninja/pgtools/connector"
    "github.com/billyninja/pgtools/scanner"
    "log"
    "time"
)


func SimpleRW(conn *connector.Connector, table string, count, rps, wps uint) {

    var r_sleep, w_sleep time.Duration
    if rps > 0 {
        r_sleep = time.Second / time.Duration(rps)
    }

    if wps > 0 {
        w_sleep = time.Second / time.Duration(wps)
    }

    params := &SimParams{
        Table:            scanner.TableName(table),
        Wipe:             WipeBefore,
        Count:            count,
        CountMode:        FillIncrement,
        ReadFunc:         ReaderGlobalCount,
        InsertsPerSecond: wps,
        ReadsPerSecond:   rps,
        SleepPerRead:     r_sleep,  // move from here, reduce friction
        SleepPerInsert:   w_sleep,  // move from here, reduce friction
    }

    Sim(conn, params)
}

func getTableByName(allTables []*scanner.Table, name scanner.TableName) *scanner.Table {
    for _, at := range allTables {
        if at.Name == name {
            return at
        }
    }
    return nil
}

func Sim(conn *connector.Connector, params *SimParams) *SimReport {
    var selectedTable *scanner.Table
    allTables := scanner.GetAllTables(conn)
    selectedTable = getTableByName(allTables, params.Table)
    if selectedTable == nil {
        log.Printf("Table specified %s not found! \n\n", params.Table)
        return nil
    }

    expected_duration := (time.Duration(params.Count) * params.SleepPerInsert)

    report := &SimReport{
        Status:             SimRunning,
        StartedAt:          time.Now(),
        SimulationParams:   params,
        UsedConnector:      conn,
        Eta:                time.Now().Add(expected_duration),
    }

    fillFKConstrains(conn, selectedTable, allTables)

    Fill(conn, selectedTable, params, report)
    fmt.Printf("%s", report)

    return report
}

func fillFKConstrains(conn *connector.Connector, tb *scanner.Table, tbs []*scanner.Table) {
    println("Solving fkdep for: ", tb.Name)
    for _, ct := range tb.Constraints {
        if ct.FTable != nil {
            fktable := getTableByName(tbs, *ct.FTable)
            fillFKConstrains(conn, fktable, tbs)
            fk_count, err := Read(conn, fktable)
            if err != nil {
                return
            }

            if fk_count < 100 {
                for i := 0; i < (100 - fk_count); i++ {
                    conn.Insert(BaseInsertQuery(fktable, 0), true)
                }
            }
        }
    }
}


