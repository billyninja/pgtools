package main

import (
	"os"
	"flag"
	"github.com/billyninja/pgtools/connector"
	"github.com/billyninja/pgtools/filler"
	"github.com/billyninja/pgtools/scanner"
	"log"
	"time"
)

var FlagDBName = flag.String("db", "", "Database to connect to")
var FlagDBHost = flag.String("host", "localhost", "Database host address")
var FlagDBUser = flag.String("user", "postgres", "Database user")
var FlagDBPass = flag.String("pass", "postgres", "Database user password")
var FlagDBPort = flag.String("port", "5432", "Database port")
var FlagSimIPS = flag.Int("Ips", 100, "Inserts Per Second")
var FlagSimRPS = flag.Int("Rps", 0, "Reads Per Second")
var FlagSimCount = flag.Int("total", 10000, "Rows to be inserted during the battery")
var FlagSimTable = flag.String("table", "", "Table to be inserted")

func Sim(conn *connector.Connector, params *filler.SimulationParams) int {
	var selectedTable *scanner.Table

	allTables := scanner.GetAllTables(conn)

	for _, at := range allTables {
		if at.Name == params.Table {
			selectedTable = at
		}
	}

	if selectedTable == nil {
		log.Printf("Table specified %s not found! \n\n", *FlagSimTable)
		return 1
	}

	filler.Fill(conn, selectedTable, params)

	return 0
}

func main() {
	flag.Parse()

	conn, _ := connector.NewConnector(
		*FlagDBHost,
		*FlagDBPort,
		*FlagDBUser,
		*FlagDBPass,
		*FlagDBName)

	if *FlagSimTable == "" {
		log.Printf("No table specified!\n\nPlease inform a table using the command-line arg -table\n\n")
		return
	}

	var rps, ips time.Duration
	if *FlagSimRPS > 0 {
		rps = time.Second / time.Duration(*FlagSimRPS)
	}

	if *FlagSimIPS > 0 {
		ips = time.Second / time.Duration(*FlagSimIPS)
	}

	sim_params := &filler.SimulationParams{
		Table:		      scanner.TableName(*FlagSimTable),
		Wipe:             filler.WipeBefore,
		Count:            *FlagSimCount,
		CountMode:        filler.FillIncrement,
		ReadFunc:         filler.ReaderGlobalCount,
		InsertsPerSecond: *FlagSimIPS,
		ReadsPerSecond:   *FlagSimRPS,
		SleepPerRead:     rps,
		SleepPerInsert:   ips,
	}
	os.Exit(Sim(conn, sim_params))
}
