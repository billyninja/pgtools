package connector

import (
	"database/sql"
	"fmt"
	"github.com/billyninja/pgtools/colors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"strings"
	"sync"
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
	FlushLock sync.Mutex
}

func (cn *Connector) CheckFlushTimeout() {
	for {
		//if !cn.FlushLock && (cn.WriteCfg.FlushTimeout > 0*time.Second && time.Since(cn.LastFlush) >= cn.WriteCfg.FlushTimeout) && len(cn.WriteAcc) > 0 {
		if (cn.WriteCfg.FlushTimeout > 0*time.Second && time.Since(cn.LastFlush) >= cn.WriteCfg.FlushTimeout) && len(cn.WriteAcc) > 0 {
			err := cn.FlushNow(true)
			if err != nil {
				log.Printf("<ERRD AT TIMEOUT ENGINE>\n\n%v\n", err)
				continue
			}
		}

		time.Sleep(1 * time.Millisecond)
	}
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

	wcfg := &WriteCfg{128, 5000 * time.Millisecond}
	cn := &Connector{
		WriteCfg: wcfg,
		WriteAcc: []string{},
		DB:       dbc,
	}
	go func() {
		cn.CheckFlushTimeout()
	}()

	return cn, nil
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

func (conn *Connector) Insert(q string, flushnow bool) (bool, bool, error) {
	persisted := false
	pos := len(conn.WriteAcc) + 1
	conn.WriteAcc = append(conn.WriteAcc, q)
	var err error
	if flushnow || (conn.WriteCfg.AccLimit > 0 && pos >= conn.WriteCfg.AccLimit) {
		err = conn.FlushNow(false)
		if err == nil {
			persisted = true
		}
	}

	return true, persisted, err
}

func (conn *Connector) DirectInsert(query string) (*sql.Rows, error) {
	rows, err := conn.DB.Query(query)
	if err != nil {
		log.Printf("ERROR: \n\n\n %s \n\n\n", query)
		return nil, err
	}

	return rows, err
}

func (conn *Connector) FlushNow(timeout bool) error {

	conn.FlushLock.Lock()
	defer conn.FlushLock.Unlock()

	trigger := "count"
	if timeout {
		trigger = "timeout"
	}

	acc := make([]string, len(conn.WriteAcc))
	copy(acc, conn.WriteAcc)
	conn.WriteAcc = []string{}

	tq := strings.Join(acc, "; ")
	t1 := time.Now()

	_, err := conn.DB.Exec(tq)
	lat := time.Since(t1)
	if err != nil {
		log.Printf("ERROR: \n\n\n %s \n\n\n", tq)

		log.Printf("Starting Drain down...\n%d queries on the queue\n\n", len(acc))
		for _, q := range acc {
			println(".")
			_, _, err = conn.Insert(q, true)
			if err != nil {
				msg := fmt.Sprintf(" Errd in Drain down: \n\n%s\n\n", err)
				log.Printf(colors.Red(msg))
			}
		}

		return err
	}
	if timeout {
		log.Printf("<%s -%s- %d - s: %d l: %s>\n",
			colors.Green("FLUSHED!"),
			colors.Yellow(trigger),
			len(acc),
			len(tq), lat)
	}

	conn.LastFlush = time.Now()
	return err
}
