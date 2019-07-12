package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/prometheus/common/log"
	"time"
)

type Config struct {
	ShowSQL      bool          `json:"show_sql"`       // show sql logs
	ShowExecTime bool          `json:"show_exec_time"` // show sql execute time
	DSN          string        `json:"dsn"`            // write data source name.
	ReadDSN      []string      `json:"read_dsn"`       // read data source name.
	Active       int           `json:"active"`         // pool
	Idle         int           `json:"idle"`           // pool
	IdleTimeout  time.Duration `json:"idle_timeout"`   // connect max life time.
	QueryTimeout time.Duration `json:"query_timeout"`  // query sql timeout
	ExecTimeout  time.Duration `json:"exec_timeout"`   // execute sql timeout
	TranTimeout  time.Duration `json:"tran_timeout"`   // transaction sql timeout
}

func NewMySQL(c *Config) *xorm.EngineGroup {
	if c.QueryTimeout == 0 || c.ExecTimeout == 0 || c.TranTimeout == 0 {
		panic("mysql must be set query/execute/transction timeout")
	}
	dataSourceSlice := []string{c.DSN}
	dataSourceSlice = append(dataSourceSlice, c.ReadDSN...)
	db, err := xorm.NewEngineGroup("mysql", dataSourceSlice)
	if err != nil {
		log.Error("open mysql error(%v)", err)
		panic(err)
	}
	db.SetMaxIdleConns(c.Idle)
	db.SetMaxOpenConns(c.Active)
	db.SetConnMaxLifetime(c.IdleTimeout)
	db.ShowSQL(c.ShowSQL)
	db.ShowExecTime(c.ShowExecTime)
	return db
}
