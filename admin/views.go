package admin

import (
    "github.com/billyninja/pgtools/scanner"
    "github.com/jmoiron/sqlx"
    "html/template"
    "io"
    "strings"
)

type EditView struct {
    Labels []template.HTML
    Inputs []template.HTML
    Table  *scanner.Table
    // Actions []SingleAction
}

type ListView struct {
    Columns    []template.HTML
    Rows       []template.HTML
    Table      *scanner.Table
    PkIndexing []int
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

func in(i int, l []int) bool {
    for _, j := range l {
        if i == j {
            return true
        }
    }
    return false
}

func mapValues(columns []*scanner.Column, pk_idx []int, rows *sqlx.Rows, rowlist *[]template.HTML, htmlgen val2html) {
    for rows.Next() {
        t_row := template.HTML("")
        sortedValues, _ := rows.SliceScan()

        pkval := ""
        for i, cl := range columns {
            fval := format_value(sortedValues[i])
            t_row += htmlgen(cl, fval)

            if len(pk_idx) > 0 {
                if in(i, pk_idx) {
                    pkval += "-" + fval
                }
            }
        }
        t_row += template.HTML(`<input type="hidden" value="` + pkval + `" />`)

        *rowlist = append(*rowlist, t_row)
    }
}

/* EDIT VIEW IMPLEMENTATION */
func NewEditView(table *scanner.Table, rows *sqlx.Rows) *EditView {

    ev := &EditView{}
    sql_cols, _ := rows.Columns()

    sortedColumns := mapCols(sql_cols, table.Columns, &ev.Labels, LabelHTML)
    mapValues(sortedColumns, []int{}, rows, &ev.Inputs, LabelAndInputHTML)

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

func buildPkIndexing(sortedColumns []*scanner.Column, pk *scanner.Constraint) []int {
    idx := []int{}
    pk_columns := strings.Split(pk.Columns, ",")

    for i, cl := range sortedColumns {
        for _, pk_cl := range pk_columns {

            if cl.Name == scanner.ColumnName(pk_cl) {
                idx = append(idx, i)
            }
        }
    }

    return idx
}

/* LIST VIEW IMPLEMENTATION */
func NewListView(table *scanner.Table, rows *sqlx.Rows) *ListView {
    lv := &ListView{}
    sql_cols, _ := rows.Columns()

    sortedColumns := mapCols(sql_cols, table.Columns, &lv.Columns, ThHTML)

    pk := table.GetPK()
    lv.PkIndexing = buildPkIndexing(sortedColumns, pk)

    mapValues(sortedColumns, lv.PkIndexing, rows, &lv.Rows, TdHTML)

    return lv
}

func (lv *ListView) GetTitle() string {
    return " " + string(lv.Table.Name) + " list view"
}

func (lv *ListView) PartialHTML(w io.Writer) error {
    var err error
    // TODO: MAKE IT EXTERNAL parameter
    // or at Creation time

    livereload := true
    if livereload {
        list_template_body, err = template.ParseFiles("/home/joao/go/src/github.com/billyninja/pgtools/admin/templates/list.html")
        if err != nil {
            return err
        }
    }

    err = list_template_body.ExecuteTemplate(w, "list.html", lv)
    if err != nil {
        println("errd ExecuteTemplate ", err)
    }
    return err
}

func (lv *ListView) CompleteHTML(buffer io.Writer) error {
    var err error
    return err
}
