package bench

import (
	"database/sql"
	"fmt"
	"github.com/billyninja/pgtools/connector"
	"github.com/billyninja/pgtools/rnd"
	"github.com/billyninja/pgtools/scanner"
	"log"
	"math/rand"
	"time"
)

var RAWREL map[scanner.TableName][]sql.RawBytes

func rndColumn(cl *scanner.Column) string {
	switch cl.Type {
	case "numeric":
		v := rnd.PSQL_numeric(99999.99, 2)
		return fmt.Sprintf("%.2f", v)
	case "integer":
		v := rnd.PSQL_int(9999999)
		return fmt.Sprintf("%d", v)
	case "character varying":
		max := 128
		if cl.CharMaxLength != nil {
			max = int(*cl.CharMaxLength)
		}
		return rnd.PSQL_var_char(2, max)
	case "date":
		return rnd.PSQL_datetime(1, 2)
	case "timestamp without time zone":
		return rnd.PSQL_datetime(0, 0)
	case "boolean":
		return rnd.PSQL_bool()
	case "json":
		return "'{}'"
	default:
		log.Panicf("\n\nRAND DATA-TYPE NOT IMPLEMENTED: %s\n\n", cl.Type)
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

		if cl.RefersTo != nil {
			n := rand.Intn(len(RAWREL[cl.RefersTo.Name]))
			values += fmt.Sprintf(`'%s'`, RAWREL[cl.RefersTo.Name][n])
		} else {
			values += rndColumn(cl)
		}

		values += ", "
	}
	base = base[0 : len(base)-2]
	values = values[0 : len(values)-2]

	base += ") "
	base += values
	base += ")"

	base += fmt.Sprintf(` RETURNING %s`, tb.PkExp)
	return base
}

func Count(conn *connector.Connector, tb *scanner.Table) (int, error) {
	sql_count := fmt.Sprintf(`SELECT COUNT(*) as cnt FROM "%s";`, tb.Name)
	rows, err := conn.Sel(sql_count)

	if err != nil {
		log.Panicf("Error at COUNT: %+v", err)
	}

	curr_cnt := &CountRS{}
	for rows.Next() {
		err := rows.StructScan(curr_cnt)
		if err != nil {
			log.Panic("err parsing table struct:\n\n ", err)
		}
	}

	return curr_cnt.Cnt, err
}

func Wipe(conn *connector.Connector, tb *scanner.Table) {
	sql_wipe := fmt.Sprintf(`DELETE FROM "%s";`, tb.Name)
	_, _, err := conn.Insert(sql_wipe, true)
	if err != nil {
		log.Panic("Error at WIPE: ", err)
	}
	conn.FlushNow(false)
}

func writeEngine(conn *connector.Connector, tb *scanner.Table, params *SimParams, report *SimReport) {

	report.writeCount = 0
	for report.writeCount < params.Count {
		t1 := time.Now()
		_, flushed, err := conn.Insert(BaseInsertQuery(tb, 0), false)
		if err != nil {
			log.Panicf("\n\n%v\n\n", err)
		}
		report.writeCount += 1

		lat := time.Since(t1)
		if report.writeCount%10 == 0 {
			smp := &Sample{
				Latency:    lat,
				WriteCount: report.writeCount,
				ReadCount:  report.readCount,
			}
			if flushed {
				report.FlushSamples = append(report.FlushSamples, smp)
			} else {
				report.InsertSamples = append(report.InsertSamples, smp)
			}
		}
		slp := (params.SleepPerInsert - lat)/2
		fmt.Printf("%s", slp)
		time.Sleep(slp)
	}
	conn.FlushNow(false)
	report.Finish()
}

func readEngine(conn *connector.Connector, tb *scanner.Table, sleep time.Duration, report *SimReport) {
	report.readCount = 0
	go func() {
		for {

			t1 := time.Now()
			_, err := Count(conn, tb)
			if err != nil {
				log.Panic("Erred at readEngine load:\n\n %v", err)
			}

			report.readCount += 1
			lat := time.Since(t1)

			if report.readCount%10 == 0 {
				report.ReadSamples = append(report.ReadSamples, &Sample{
					Latency:    lat,
					WriteCount: report.writeCount,
					ReadCount:  report.readCount,
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

	report.StartedAt = time.Now()
	if params.InsertsPerSecond > 0 {
		writeEngine(conn, tb, params, report)
	}

	if params.Wipe == WipeAfter || params.Wipe == WipeBeforeAndAfter {
		Wipe(conn, tb)
	}
}

func fillFKConstrains(conn *connector.Connector, tb *scanner.Table, tbs []*scanner.Table) {
	rand.Seed(time.Now().UnixNano())

	println("Filling fk Constraints for ", tb.Name)

	for _, ct := range tb.Constraints {
		if ct.FTable != nil {
			fktable := getTableByName(tbs, *ct.FTable)
			fillFKConstrains(conn, fktable, tbs)
			fk_count, err := Count(conn, fktable)
			if err != nil {
				return
			}

			if fk_count < 50 {
				for i := 0; i < (50 - fk_count); i++ {
					qry := BaseInsertQuery(fktable, 0)
					rows, err := conn.DirectInsert(qry)
					if err != nil {
						log.Panic("Errd at fillFKConstrains - DirectInsert - ", err)
					}

					for rows.Next() {
						var raw sql.RawBytes
						err = rows.Scan(&raw)
						if err != nil {
							log.Panic("Errd at fillFKConstrains - DirectInsert - ", err)
						}

						key := *ct.FTable
						if _, ok := RAWREL[key]; ok {
							RAWREL[key] = append(RAWREL[key], raw)
						} else {
							RAWREL[key] = []sql.RawBytes{raw}
						}
					}
				}
			}
		}
	}
}
