package main

import (
    "fmt"
    "log"
    "flag"
//    "strings"
    _ "github.com/lib/pq"
    "github.com/jmoiron/sqlx"
)


type Column struct {
    Name        string      `db:"column_name"`
    Type        string      `db:"data_type"`
    Default     *string     `db:"column_default"`
    Nullable    string      `db:"is_nullable"`
}

type Table struct {
    Name        string      `db:"table_name"`
    TableType   string      `db:"table_type"`
    Columns     []*Column
}

func (tb *Table) BaseInsertQuery() string {
    base := fmt.Sprintf(`INSERT INTO "%s" (`, tb.Name)
    nc := len(tb.Columns)
    for i, c := range tb.Columns {
        base += c.Name
        if nc > i + 1  {
            base += ", "
        }
    }
    base += ") VALUES (...)"

    return base
}

func (t Table) String() string {
    out := fmt.Sprintf("Table: %s\n\tColumns:\n", t.Name)
    for _, c := range t.Columns {
        dft := "[no default]"
        if c.Default != nil {
            dft = *c.Default
        }
        out += fmt.Sprintf("\t\t%s (%s) dft: %s  Null? %s\n", c.Name, c.Type, dft, c.Nullable)
    }
    return out
}

var FlagDBName = flag.String("db", "", "Database to connect to")
var FlagDBHost = flag.String("host", "localhost", "Database host address")
var FlagDBUser = flag.String("user", "postgres", "Database user")
var FlagDBPass  = flag.String("pass", "postgres", "Database user password")
var FlagDBPort = flag.String("port", "5432", "Database port")

func sel(db *sqlx.DB, q string) (*sqlx.Rows, error) {
    var rows *sqlx.Rows
    rows, err := db.Queryx(q)
    if err != nil {
        log.Println("%v", err)
        log.Println(q)
    }

    return rows, err
}


func main() {
    flag.Parse()

    strConn := fmt.Sprintf(
        "dbname=%s user=%s password=%s host=%s port=%s sslmode=disable",
        *FlagDBName,
        *FlagDBUser,
        *FlagDBPass,
        *FlagDBHost,
        *FlagDBPort,
    )

    db := sqlx.MustOpen("postgres", strConn)
    qTables := "SELECT table_name, table_type FROM information_schema.tables WHERE table_schema='public';"
    rows, err := sel(db, qTables)
    if err != nil {
        log.Panic("-")
    }

    var allTables []*Table

    for rows.Next() {
        tb := &Table{}
        err := rows.StructScan(tb)
        if err != nil {
            log.Panic("---\n%v\n----", err)
        }
        allTables = append(allTables, tb)
    }

    for _, tb := range allTables {
        qColumns := fmt.Sprintf(`
            SELECT column_name, data_type, column_default, is_nullable FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = '%s'
        `, tb.Name)

        sel(db, qColumns)
        rows, err := sel(db, qColumns)
        if err != nil {
            log.Panic("err parsing table struct:\n %v", err)
        }

        for rows.Next() {
            cl := &Column{}
            err := rows.StructScan(cl)
            if err != nil {
                log.Panic("err parsing table struct:\n %v", err)
            }
            tb.Columns = append(tb.Columns, cl)
        }
    }
    fmt.Printf("%s\n", allTables[0])
    fmt.Printf("%s\n", allTables[0].BaseInsertQuery())
}
