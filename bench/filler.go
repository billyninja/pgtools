package bench

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
		base += `"` + string(cl.Name) + `"`
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

type Count struct {
	Cnt int `db:"cnt"`
}

func Read(conn *connector.Connector, tb *scanner.Table) {
	sql_count := fmt.Sprintf(`SELECT COUNT(*) as cnt FROM "%s";`, tb.Name)
	rows, err := conn.Sel(sql_count)

	if err != nil {
		log.Panic("Error at COUNT: %+v", err)
	}

	for rows.Next() {
		curr_cnt := &Count{}
		err := rows.StructScan(curr_cnt)
		if err != nil {
			log.Panic("err parsing table struct:\n %v", err)
		}
	}
}

func Wipe(conn *connector.Connector, tb *scanner.Table) {
	sql_wipe := fmt.Sprintf(`DELETE FROM "%s";`, tb.Name)
	_, _, err := conn.Insert(sql_wipe)
	if err != nil {
		log.Panic("Error at WIPE: %+v", err)
	}
	conn.FlushNow(false)
}

func writeEngine(conn *connector.Connector, tb *scanner.Table, params *SimParams, report *SimReport) {

	report.writeCount = 0
	for report.writeCount < params.Count {
		t1 := time.Now()
		_, flushed, err := conn.Insert(BaseInsertQuery(tb, 0))
		if err != nil {
			log.Panic("\n\n%v\n\n", err)
		}
		report.writeCount += 1
		lat := time.Since(t1)

		if report.writeCount%10 == 0 {
			smp := &Sample{
				Latency: lat,
				WriteCount: report.writeCount,
				ReadCount: report.readCount,
			}
			if flushed {
				report.FlushSamples = append(report.FlushSamples, smp)
			} else {
				report.InsertSamples = append(report.InsertSamples, smp)
			}
		}

		time.Sleep(params.SleepPerInsert - lat)
	}
	conn.FlushNow(false)
	report.Finish()
}

func readEngine(conn *connector.Connector, tb *scanner.Table, sleep time.Duration, report *SimReport) {
	report.readCount = 0
	go func() {
		for {

			t1 := time.Now()
			Read(conn, tb)
			report.readCount += 1
			lat := time.Since(t1)

			if report.readCount%10 == 0 {
				report.ReadSamples = append(report.ReadSamples, &Sample{
					Latency: lat,
					WriteCount: report.writeCount,
					ReadCount: report.readCount,
				})
			}

			time.Sleep(sleep - lat)
		}
	}()
}

func Fill(conn *connector.Connector, tb *scanner.Table, params *SimParams, report *SimReport) {
	rand.Seed(time.Now().UnixNano())

	if params.Wipe == WipeBefore || params.Wipe == WipeBeforeAndAfter {
		Wipe(conn, tb)
	}

	if params.ReadsPerSecond > 0 {
		readEngine(conn, tb, params.SleepPerRead, report)
	}

	if params.InsertsPerSecond > 0 {
		writeEngine(conn, tb, params, report)
	}

	if params.Wipe == WipeAfter || params.Wipe == WipeBeforeAndAfter {
		Wipe(conn, tb)
	}
}
