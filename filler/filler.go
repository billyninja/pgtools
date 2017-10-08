package filler

import (
	"fmt"
    "math/rand"
	"github.com/billyninja/pgtools/connector"
	"github.com/billyninja/pgtools/rnd"
	"github.com/billyninja/pgtools/scanner"
	"log"
	"time"
)

func BaseInsertQuery(tb *scanner.Table) string {
	base := fmt.Sprintf(`INSERT INTO "%s" (`, tb.Name)
	nc := len(tb.Columns)
	for i, c := range tb.Columns {
		base += c.Name
		if nc > i+1 {
			base += ", "
		}
	}

    v1 := fmt.Sprintf(
        "(%s, %s, %s, %s, %s, %s, %s, %.2f, %.2f, %.2f , '{}', FALSE)",
        rnd.PSQL_var_char(3, 3),
        rnd.PSQL_var_char(3, 3),
        rnd.PSQL_datetime(1, 2),
        rnd.PSQL_datetime(1, 2),
        rnd.PSQL_datetime(0, 0),
        rnd.PSQL_datetime(0, 0),
        rnd.PSQL_var_char(2, 2),
        rnd.PSQL_numeric(99999.99, 2),
        rnd.PSQL_numeric(99.99, 2),
        rnd.PSQL_numeric(99.99, 2),
    )

	base += ") VALUES "
    base += v1

	return base
}

func Fill(conn *connector.Connector, tb *scanner.Table, nrows int64) {
	rand.Seed(time.Now().UnixNano())
    i := int64(0)
	for i < nrows {
		_, _, err := conn.Insert(BaseInsertQuery(tb))
		if err != nil {
			log.Printf("\n\n\n\n%v\n\n\n\n", err)
			return
		}
		time.Sleep(1 * time.Millisecond)
		i += 1
	}
    conn.FlushNow()
}
