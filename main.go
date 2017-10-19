package main

import (
	"flag"
	"github.com/billyninja/pgtools/connector"
	"github.com/billyninja/pgtools/filler"
	"github.com/billyninja/pgtools/scanner"
	"time"
)

var FlagDBName = flag.String("db", "", "Database to connect to")
var FlagDBHost = flag.String("host", "localhost", "Database host address")
var FlagDBUser = flag.String("user", "postgres", "Database user")
var FlagDBPass = flag.String("pass", "postgres", "Database user password")
var FlagDBPort = flag.String("port", "5432", "Database port")
var FlagSimIPS = flag.Int("Ips", 100, "Inserts Per Second")
var FlagSimRPS = flag.Int("Rps", 100, "Reads Per Second")
var FlagSimCount = flag.Int("total", 10000, "Rows to be inserted during the battery")

func main() {
	flag.Parse()

	conn, _ := connector.NewConnector(
		*FlagDBHost,
		*FlagDBPort,
		*FlagDBUser,
		*FlagDBPass,
		*FlagDBName)

	allTables := scanner.GetAllTables(conn)

	sim_params := &filler.SimulationParams{
		Wipe:           	filler.WipeBefore,
		Count:		    	2000,
		CountMode:      	filler.FillIncrement,
		InsertsPerSecond: 	*FlagSimIPS,
		ReadsPerSecond: 	*FlagSimRPS,
		SleepPerRead: 		time.Second/time.Duration(*FlagSimRPS),
		SleepPerInsert: 	time.Second/time.Duration(*FlagSimIPS),
	}

	filler.Fill(conn, allTables[0], sim_params)
}
