package filler

import (
	"fmt"
	"github.com/billyninja/pgtools/connector"
	"github.com/billyninja/pgtools/rnd"
	"github.com/billyninja/pgtools/scanner"
	"log"
	"math/rand"
	"time"
)

func rndColumn(cl *scanner.Column) string {
	switch cl.Type {
	case "numeric":
		v := rnd.PSQL_numeric(99999.99, 2)
		return fmt.Sprintf("%.2f", v)
	case "character varying":
		max := 128
		if cl.CharMaxLength != nil {
			max = int(*cl.CharMaxLength)
		}
		return rnd.PSQL_var_char(1, max)
	case "date":
		return rnd.PSQL_datetime(1, 2)
	case "timestamp without time zone":
		return rnd.PSQL_datetime(0, 0)
	case "boolean":
		return rnd.PSQL_bool()
	case "json":
		return "'{}'"
	default:
		fmt.Printf("\n\nRAND DATA-TYPE NOT IMPLEMENTED: %s\n\n", cl.Type)
		break
	}

	return ""
}

func BaseInsertQuery(tb *scanner.Table, skip_nullable uint8) string {
	base := fmt.Sprintf(`INSERT INTO "%s" (`, tb.Name)
	values := "VALUES ("
	for _, cl := range tb.Columns {
		if skip_nullable > 0 && cl.Nullable == "YES" {
			continue
		}
		base += cl.Name
		base += ", "
		values += rndColumn(cl)
		values += ", "
	}
	base = base[0 : len(base)-2]
	values = values[0 : len(values)-2]

	base += ") "
	base += values
	base += ")"

	return base
}

func Fill(conn *connector.Connector, tb *scanner.Table, nrows int64) {
	rand.Seed(time.Now().UnixNano())
	i := int64(0)
	for i < nrows {
		_, _, err := conn.Insert(BaseInsertQuery(tb, 0))
		if err != nil {
			log.Printf("\n\n\n\n%v\n\n\n\n", err)
			return
		}
		time.Sleep(100 * time.Millisecond)
		i += 1
	}
	conn.FlushNow()
}
