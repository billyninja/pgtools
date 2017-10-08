package filler

import (
    "fmt"
    "log"
    "time"
    "github.com/billyninja/pgtools/connector"
    "github.com/billyninja/pgtools/scanner"
    "github.com/billyninja/pgtools/rnd"
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
    base += ") VALUES ('BSB', 'NAT', '2017-11-23', '2018-01-21', NOW(), NOW(), 'JJ', 1030.23, 27.79, 103.02 , '{}', FALSE)"

    "('BSB', 'NAT', '2017-11-23', '2018-01-21', NOW(), NOW(), 'JJ', 1030.23, 27.79, 103.02 , '{}', FALSE)"
    PSQL_var_char(3),
    PSQL_var_char(3),
    PSQL_datetime(1, 2),
    PSQL_datetime(0, 0),
    PSQL_var_char(2, 2),
    PSQL_numeric(99999.99, 2),
    PSQL_numeric(99.99, 2),
    PSQL_numeric(999.99, 2),
    "{}", // --- priceline static

    return base
}

func Fill(conn *connector.Connector , tb *scanner.Table, nrows int64) {
    i := int64(0)
    for i < nrows {
        _, _, err := conn.Insert(BaseInsertQuery(tb))
        if err != nil {
            log.Printf("\n\n\n\n%v\n\n\n\n", err)
            return
        }
        time.Sleep(120 * time.Millisecond)
        i += 1
    }
}
