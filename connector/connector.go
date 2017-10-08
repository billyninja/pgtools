package connector

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"strings"
	"time"
)

type WriteCfg struct {
	AccLimit     int
	FlushTimeout time.Duration
}

type Connector struct {
	DB        *sqlx.DB
	WriteCfg  *WriteCfg
	WriteAcc  []string
	LastFlush time.Time
}

func NewConnector(host, port, user, pass, db string) (*Connector, error) {
	strConn := fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%s sslmode=disable",
		db,
		user,
		pass,
		host,
		port,
	)
	dbc := sqlx.MustOpen("postgres", strConn)

	wcfg := &WriteCfg{100, 10 * time.Second}
	return &Connector{
		WriteCfg: wcfg,
		WriteAcc: []string{},
		DB:       dbc,
	}, nil
}

func (conn *Connector) Sel(q string) (*sqlx.Rows, error) {
	var rows *sqlx.Rows
	rows, err := conn.DB.Queryx(q)
	if err != nil {
		log.Println("%v", err)
		log.Println(q)
	}

	return rows, err
}

func (conn *Connector) Insert(q string) (bool, bool, error) {
	persisted := false
	pos := len(conn.WriteAcc) + 1
	if (conn.WriteCfg.AccLimit > 0 && pos >= conn.WriteCfg.AccLimit) || (conn.WriteCfg.FlushTimeout > time.Second*0 && time.Since(conn.LastFlush) >= conn.WriteCfg.FlushTimeout) {
		tq := strings.Join(conn.WriteAcc, "; ")
		t1 := time.Now()
		_, err := conn.DB.Exec(tq)
		lat := time.Since(t1)
		if err != nil {
			log.Println(tq)
			return false, false, err
		}
		persisted = true
		conn.WriteAcc = []string{}
		conn.LastFlush = time.Now()
		log.Printf("<PERSISTED! %d - s: %d l: %s>\n", pos, len(tq), lat)
	} else {
		conn.WriteAcc = append(conn.WriteAcc, q)
		log.Println("<ACCD!>")
	}

	return true, persisted, nil
}
