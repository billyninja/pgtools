package admin

import (
    //    "fmt"
    //    "time"
    "github.com/billyninja/pgtools/scanner"
    "github.com/jmoiron/sqlx"
    "html/template"
    "io"
)

type EditView struct {
    Labels []template.HTML
    Inputs []template.HTML
    Table  *scanner.Table
    // Actions []SingleAction
}

type ListView struct {
    Columns []template.HTML
    Rows    []template.HTML
    Table   *scanner.Table
    // Actions []BulkAction
}

type SingleView struct {
    Labels []template.HTML
    Values []template.HTML
    Table  *scanner.Table
    // Actions []SingleAction
}

func mapCols(sql_cols []string, columns []*scanner.Column, mylist *[]template.HTML, htmlgen col2html) []*scanner.Column {
    sortedColumns := make([]*scanner.Column, len(sql_cols))

    for i, scl := range sql_cols {
        for _, cl := range columns {
            if scl == string(cl.Name) {
                sortedColumns[i] = cl
                *mylist = append(*mylist, htmlgen(cl))
            }
        }
    }

    return sortedColumns
}

func mapValues(columns []*scanner.Column, rows *sqlx.Rows, mylist *[]template.HTML, htmlgen val2html) {
    for rows.Next() {
        t_row := template.HTML("")
        sortedValues, _ := rows.SliceScan()
        for i, cl := range columns {
            t_row += htmlgen(cl, sortedValues[i])
        }
        *mylist = append(*mylist, t_row)
    }
}

/* EDIT VIEW IMPLEMENTATION */
func NewEditView(table *scanner.Table, rows *sqlx.Rows) *EditView {

    ev := &EditView{}
    sql_cols, _ := rows.Columns()

    sortedColumns := mapCols(sql_cols, table.Columns, &ev.Labels, LabelHTML)
    mapValues(sortedColumns, rows, &ev.Inputs, LabelAndInputHTML)

    return ev
}

func (ev *EditView) GetTitle() string {
    return " " + string(ev.Table.Name) + " edit view"
}

func (ev *EditView) PartialHTML(buffer io.Writer) error {
    var err error

    // TODO: MAKE IT EXTERNAL parameter
    // or at Creation time

    livereload := true
    if livereload {
        edit_template_body, err = template.ParseFiles("admin/templates/edit.html")
        if err != nil {
            return err
        }
    }
    err = edit_template_body.ExecuteTemplate(buffer, "edit.html", ev)
    if err != nil {
        println("errd ExecuteTemplate ", err)
    }

    return err
}

func (ev *EditView) CompleteHTML(buffer io.Writer) {

}

/* LIST VIEW IMPLEMENTATION */
func NewListView(table *scanner.Table, rows *sqlx.Rows) *ListView {
    lv := &ListView{}
    sql_cols, _ := rows.Columns()

    sortedColumns := mapCols(sql_cols, table.Columns, &lv.Columns, ThHTML)
    mapValues(sortedColumns, rows, &lv.Rows, TdHTML)

    return lv
}

func (lv *ListView) GetTitle() string {
    return " " + string(lv.Table.Name) + " list view"
}

func (lv *ListView) PartialHTML(buffer io.Writer) error {
    var err error

    // TODO: MAKE IT EXTERNAL parameter
    // or at Creation time

    livereload := true
    if livereload {
        list_template_body, err = template.ParseFiles("admin/templates/list.html")
        if err != nil {
            return err
        }
    }
    err = list_template_body.ExecuteTemplate(buffer, "list.html", lv)
    if err != nil {
        println("errd ExecuteTemplate ", err)
    }

    return err
}

func (lv *ListView) CompleteHTML(buffer io.Writer) error {
    var err error
    return err
}
