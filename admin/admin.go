package admin

import (
    "bytes"
    "fmt"
    "github.com/billyninja/pgtools/connector"
    "github.com/billyninja/pgtools/scanner"
    "github.com/jmoiron/sqlx"
    "html/template"
    "log"
)

func ScratchAdmin(conn *connector.Connector) {
    allTables := scanner.GetAllTables(conn)

    for _, tb := range allTables {
        q := QueryListAll(tb.Name)
        rows, err := conn.Sel(q)
        if err != nil {
            log.Printf("Couldn't query table data!")
        }

        // ListEntries(tb, rows)
        EditEntry(tb, rows)
    }
}

func QueryListAll(table scanner.TableName) string {
    return fmt.Sprintf(`SELECT * FROM %s;`, table)
}

var edit_template_body, list_template_body *template.Template

func EditEntry(table *scanner.Table, rows *sqlx.Rows) error {
    var err error

    buffer := &bytes.Buffer{}
    ev := NewEditView(table, rows)
    ev.PartialHTML(buffer)
    fmt.Printf(">%s", buffer)

    return err
}

func ViewEntry(sortedColumns []*scanner.Column, sortedValues []interface{}) {}
