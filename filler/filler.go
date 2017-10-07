package filler

import (
    "fmt"
    "github.com/billyninja/pgtools/scanner"
)

func BaseInsertQuery(tb *scanner.Table) string {
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

func Fill(tb *scanner.Table, nrows int64) {
    i := int64(0)
    for i < nrows {
        println(BaseInsertQuery(tb))
        i += 1
    }
}
