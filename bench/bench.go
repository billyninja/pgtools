package bench

import (
	"database/sql"
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
		SleepPerRead:     r_sleep, // move from here, reduce friction
		SleepPerInsert:   w_sleep, // move from here, reduce friction
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
		Status:           SimRunning,
		SimulationParams: params,
		UsedConnector:    conn,
		Eta:              time.Now().Add(expected_duration),
	}

	RAWREL = make(map[scanner.TableName][]sql.RawBytes)
	fillFKConstrains(conn, selectedTable, allTables)

	Fill(conn, selectedTable, params, report)
	fmt.Printf("%s", report)

	return report
}

func ScratchAdmin(conn *connector.Connector) {
	allTables := scanner.GetAllTables(conn)
	actions := []string{"List", "View", "Create", "Edit", "Delete"}

	for _, tb := range allTables {

		println(">Model: ", tb.Name)
		for _, act := range actions {
			println("\tAction:", act)
			for _, cl := range tb.Columns {
				if act == "Create" || act == "Edit" {
					println("\t\t", cl.Input(""))
				} else if act == "View" {
					println("\t\tDisplay: ", cl.Name)
				} else if act == "List" {
					ListEntries(conn, tb)
					//println("\t\tList column: ", cl.Name)
				}
			}
			println("----")
		}
	}
}

func QueryListAll(table scanner.TableName) string {
	return fmt.Sprintf(`SELECT * FROM %s`, table)
}

func ListEntries(conn *connector.Connector, table *scanner.Table) {
	q := QueryListAll(table.Name)
	rows, err := conn.Sel(q)
	if err != nil {
		log.Printf("Couldn't query table data!")
	}

	sortedColumns := []*scanner.Column{}
	cols, _ := rows.Columns()
	for _, Xcl := range cols {
		for _, Tcl := range table.Columns {
			if Xcl == string(Tcl.Name) {
				sortedColumns = append(sortedColumns, Tcl)
			}
		}
	}

	for rows.Next() {

		values, _ := rows.SliceScan()
		EditEntry(sortedColumns, values)

		// for idx, val := range values {
		// 	cl := sortedColumns[idx]
		// 	if cl != nil {
		// 		fmt.Printf("%s: %s\n", cl.Name, val)
		// 	}
		// }
	}
}

func EditEntry(sortedColumns []*scanner.Column, sortedValues []interface{}) {
	println("begin edit===")
	for i, cl := range sortedColumns {
		println(cl.Input(sortedValues[i]))
	}
	println("end edit===")
}
