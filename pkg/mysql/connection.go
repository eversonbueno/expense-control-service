package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLInterface interface {
	Connect() (*sql.DB, error)
	Close() error
}

type MySQL struct {
	host         string
	user         string
	password     string
	databaseName string
	maxOpenConns int
	maxIdleConns int
	maxLifetime  time.Duration
	conn         *sql.DB
}

func New(host, user, password, dbname string, maxOpenConns, maxIdleConns int, maxLifetime time.Duration) MySQLInterface {
	return &MySQL{
		host:         host,
		user:         user,
		password:     password,
		databaseName: dbname,
		maxOpenConns: maxOpenConns,
		maxIdleConns: maxIdleConns,
		maxLifetime:  maxLifetime,
	}
}

func (m *MySQL) Connect() (*sql.DB, error) {
	if m.conn != nil {
		return m.conn, fmt.Errorf("connection already established")
	}

	connectionString := m.buildConnectionString()
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	db.SetMaxOpenConns(m.maxOpenConns)
	db.SetMaxIdleConns(m.maxIdleConns)
	db.SetConnMaxLifetime(m.maxLifetime)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	m.conn = db
	return m.conn, nil
}

func (m *MySQL) buildConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		m.user, m.password, m.host, m.databaseName)
}

func (m *MySQL) Close() error {
	if m.conn == nil {
		return fmt.Errorf("no active connection to close")
	}
	err := m.conn.Close()
	m.conn = nil
	return err
}

