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
var FlagSimSleep = flag.Int("sleep", 100, "Sleep timeout")

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
		Wipe:           filler.WipeBefore,
		Count:		    2000,
		CountMode:      filler.FillIncrement,
		SleepPerInsert: time.Millisecond * time.Duration(*FlagSimSleep),
	}

	filler.Fill(conn, allTables[0], sim_params)
}
