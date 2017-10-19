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


type WipeMode uint8
type CountMode uint8

const (
	WipeNever 			WipeMode = iota
	WipeBefore
	WipeAfter
	WipeBeforeAndAfter

	FillIncrement		CountMode = iota
	FillUntil
)


type SimulationParams struct {
	Wipe 				WipeMode
	Count 				int64
	CountMode   		CountMode
	InsertsPerSecond	int
	ReadsPerSecond		int
	SleepPerInsert		time.Duration
	SleepPerRead 		time.Duration
}


type FillReport struct {
	StartTime 		time.Time
	EndTime 		time.Time
	TotalWrites		uint64
    TotalReads		uint64
}


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


type Count struct {
	Cnt int `db:"cnt"`
}

func Read(conn *connector.Connector, tb *scanner.Table) {
	sql_count := fmt.Sprintf(`SELECT COUNT(*) as cnt FROM "%s";`, tb.Name)
	t1 := time.Now()
	rows, err := conn.Sel(sql_count)
	latency := time.Since(t1)

	if err != nil {
		log.Panic("Error at COUNT: %+v", err)
	}

	for rows.Next() {
		curr_cnt := &Count{}
		err := rows.StructScan(curr_cnt)
		if err != nil {
			log.Panic("err parsing table struct:\n %v", err)
		}
		log.Printf("Count >> %d %v", curr_cnt.Cnt, latency)
	}
}

func Wipe(conn *connector.Connector, tb *scanner.Table) {
	sql_wipe := fmt.Sprintf(`DELETE FROM "%s";`, tb.Name)
	_, _, err := conn.Insert(sql_wipe)
	if err != nil {
		log.Panic("Error at WIPE: %+v", err)
	}
	conn.FlushNow()
}

func ReadEngine(conn *connector.Connector, tb *scanner.Table, sleep time.Duration) {
	go func(){
		for {
			Read(conn, tb)
			time.Sleep(sleep)
		}
	}()
}


func Fill(conn *connector.Connector, tb *scanner.Table, params *SimulationParams) {
	rand.Seed(time.Now().UnixNano())
	Read(conn, tb)

	if params.Wipe == WipeBefore || params.Wipe == WipeBeforeAndAfter {
		Wipe(conn, tb)
	}

	ReadEngine(conn, tb, params.SleepPerRead)

	i := int64(0)
	for i < params.Count {
		_, _, err := conn.Insert(BaseInsertQuery(tb, 0))
		if err != nil {
			log.Panic("\n\n\n\n%v\n\n\n\n", err)
		}
		time.Sleep(params.SleepPerInsert)
		i += 1
	}
	conn.FlushNow()

	if params.Wipe == WipeAfter || params.Wipe == WipeBeforeAndAfter {
		Wipe(conn, tb)
	}
}
