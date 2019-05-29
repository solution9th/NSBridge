package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/solution9th/NSBridge/service/mysql"
)

var (
	ErrTxHasBegan = errors.New("transaction already begin")
	ErrTxDone     = errors.New("transaction not begin")
)

const (
	tableDomain = "dns_domain"
	tableRecord = "dns_record"
	tableAuth   = "auth"
)

// db querier
type dbQuerier interface {
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type txer interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type txEnder interface {
	Commit() error
	Rollback() error
}

type Tables struct {
	DB   dbQuerier
	isTx bool
}

func New() *Tables {
	return &Tables{
		DB: mysql.DefaultDB,
	}
}

// Begin transaction
func (t *Tables) Begin() error {

	if t.isTx {
		return ErrTxHasBegan
	}

	var tx *sql.Tx
	tx, err := t.DB.(txer).Begin()
	if err != nil {
		return err
	}
	t.isTx = true

	t.DB = tx

	return nil
}

// Commit transaction
func (t *Tables) Commit() error {
	if !t.isTx {
		return ErrTxDone
	}
	err := t.DB.(txEnder).Commit()
	if err == nil {
		t.isTx = false
		t.DB = mysql.DefaultDB
	} else if err == sql.ErrTxDone {
		return ErrTxDone
	}
	return err
}

// Rollback transaction
func (t *Tables) Rollback() error {
	if !t.isTx {
		return ErrTxDone
	}
	err := t.DB.(txEnder).Rollback()
	if err == nil {
		t.isTx = false
		t.DB = mysql.DefaultDB
	} else if err == sql.ErrTxDone {
		return ErrTxDone
	}
	return err
}
