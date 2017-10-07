package main

import (
    "flag"
    "github.com/billyninja/pgtools/connector"
    "github.com/billyninja/pgtools/scanner"
    "github.com/billyninja/pgtools/filler"
)

var FlagDBName = flag.String("db", "", "Database to connect to")
var FlagDBHost = flag.String("host", "localhost", "Database host address")
var FlagDBUser = flag.String("user", "postgres", "Database user")
var FlagDBPass  = flag.String("pass", "postgres", "Database user password")
var FlagDBPort = flag.String("port", "5432", "Database port")


func main() {
    flag.Parse()

    conn, _ := connector.NewConnector(
        *FlagDBHost,
        *FlagDBPort,
        *FlagDBUser,
        *FlagDBPass,
        *FlagDBName)

    allTables := scanner.GetAllTables(conn)
    filler.Fill(allTables[0], 2000)
}
