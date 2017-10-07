package connector

import (
    "fmt"
    "time"
    _ "github.com/lib/pq"
    "github.com/jmoiron/sqlx"
)

type WriteCfg struct {
    AccLimit        int
    FlushTimeout    time.Duration
}

type Connector struct {
    WriteCfg    *WriteCfg
    WriteAcc    []string
    DB      *sqlx.DB
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

    return &Connector{
        WriteCfg: &WriteCfg{100, 10 * time.Second},
        DB: dbc,
    }, nil
}
