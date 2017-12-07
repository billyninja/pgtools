package admin

import (
    "bytes"
    "fmt"
    "github.com/billyninja/pgtools/connector"
    "github.com/billyninja/pgtools/scanner"
    "html/template"
    "log"
)

func ScratchAdmin(conn *connector.Connector) {
    allTables := scanner.GetAllTables(conn)

    // actions := []string{"List", "View", "Create", "Edit", "Delete"}

    for _, tb := range allTables {
        if tb.Name == "compra" {
            ListEntries(conn, tb)
        }
    }
}

func QueryListAll(table scanner.TableName) string {
    return fmt.Sprintf(`SELECT * FROM %s ORDER BY data DESC LIMIT 10000;`, table)
}

var edit_template_body, list_template_body *template.Template

func ListEntries(conn *connector.Connector, table *scanner.Table) error {
    buffer := &bytes.Buffer{}

    q := QueryListAll(table.Name)
    rows, err := conn.Sel(q)
    if err != nil {
        log.Printf("Couldn't query table data!")
    }

    lv := NewListView(table, rows)
    lv.PartialHTML(buffer)
    fmt.Printf("%s", buffer)

    // sortedColumns := []*scanner.Column{}
    // cols, _ := rows.Columns()
    // for _, Xcl := range cols {
    //     for _, Tcl := range table.Columns {
    //         if Xcl == string(Tcl.Name) {
    //             sortedColumns = append(sortedColumns, Tcl)
    //             list_view.Headers = append(list_view.Headers, Tcl.ThHTML())
    //         }
    //     }
    // }

    // livereload := true
    // if livereload {
    //     list_template_body, err = template.ParseFiles("admin/templates/list.html")
    //     if err != nil {
    //         log.Printf("errd parsing template file! %s", err)
    //         return err
    //     }
    // }

    // for rows.Next() {
    //     t_row := template.HTML("")
    //     sortedValues, _ := rows.SliceScan()
    //     for i, cl := range sortedColumns {
    //         t_row = template.HTML(fmt.Sprintf("%s%s", t_row, cl.TdHTML(sortedValues[i])))
    //     }
    //     list_view.Rows = append(list_view.Rows, t_row)
    // }

    // err = list_template_body.ExecuteTemplate(buffer, "list.html", list_view)
    // fmt.Printf("%s", buffer)

    return err
}

func EditEntry(sortedColumns []*scanner.Column, sortedValues []interface{}) error {
    var err error

    // livereload := true
    // if livereload {
    //     edit_template_body, err = template.ParseFiles("admin/templates/edit.html")
    //     if err != nil {
    //         log.Printf("errd parsing template file! %s", err)
    //         return err
    //     }
    // }

    // buffer := &bytes.Buffer{}
    // inputs := []template.HTML{}
    // for i, cl := range sortedColumns {
    //     inputs = append(inputs, cl.InputHTML(sortedValues[i]))
    // }
    // err = edit_template_body.ExecuteTemplate(buffer, "edit.html", inputs)

    return err
}

func ViewEntry(sortedColumns []*scanner.Column, sortedValues []interface{}) {}
