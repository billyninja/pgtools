package admin

import (
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

/* EDIT VIEW IMPLEMENTATION */
func NewEditView(table *scanner.Table, rows sqlx.Rows) *EditView {
    return &EditView{}
}

func (ev *EditView) GetTitle() string {
    return " " + string(ev.Table.Name) + " edit view"
}

func (ev *EditView) PartialHTML(buffer io.Writer) {

}

func (ev *EditView) CompleteHTML(buffer io.Writer) {

}

func mapCols(sql_cols []string, columns []*scanner.Column, mylist []template.HTML, htmlgen col2html) []*scanner.Column {
    sortedColumns := make([]*scanner.Column, len(sql_cols))
    mylist = make([]template.HTML, len(sql_cols))
    for i, scl := range sql_cols {
        for _, cl := range columns {
            if scl == string(cl.Name) {
                sortedColumns[i] = cl
                mylist[i] = htmlgen(cl)
            }
        }
    }

    return sortedColumns
}

func mapValues(columns []*scanner.Column, rows *sqlx.Rows, mylist []template.HTML, htmlgen val2html) {
    for rows.Next() {
        t_row := template.HTML("")
        sortedValues, _ := rows.SliceScan()
        for i, cl := range columns {
            t_row += htmlgen(cl, sortedValues[i])
        }
        mylist = append(mylist, t_row)
    }
}

/* LIST VIEW IMPLEMENTATION */
func NewListView(table *scanner.Table, rows *sqlx.Rows) *ListView {
    lv := &ListView{}
    sql_cols, _ := rows.Columns()

    sortedColumns := mapCols(sql_cols, table.Columns, lv.Columns, ThHTML)
    mapValues(sortedColumns, rows, lv.Rows, TdHTML)

    return lv
}

func (lv *ListView) GetTitle() string {
    return " " + string(lv.Table.Name) + " list view"
}

func (lv *ListView) PartialHTML(buffer io.Writer) error {
    var err error
    return err
}

func (lv *ListView) CompleteHTML(buffer io.Writer) error {
    var err error
    return err
}
