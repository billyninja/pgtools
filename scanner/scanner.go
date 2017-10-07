package scanner

import (
    "fmt"
    "log"
    "github.com/billyninja/pgtools/connector"
    "github.com/jmoiron/sqlx"
)


func sel(conn *connector.Connector, q string) (*sqlx.Rows, error) {
    var rows *sqlx.Rows
    rows, err := conn.DB.Queryx(q)
    if err != nil {
        log.Println("%v", err)
        log.Println(q)
    }

    return rows, err
}

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

func GetAllTables(conn *connector.Connector) []*Table {
    qTables := "SELECT table_name, table_type FROM information_schema.tables WHERE table_schema='public';"
    rows, err := sel(conn, qTables)
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

        rows, err := sel(conn, qColumns)
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

    return allTables
}
