package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Database SQLite数据库连接包装器
type Database struct {
	db *sql.DB
}

// New 创建新的SQLite数据库连接
func New(databasePath string) (*Database, error) {
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(1) // SQLite 只支持单个写连接
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// 启用 WAL 模式和外键约束
	if err := enablePragmas(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set pragmas: %w", err)
	}

	return &Database{db: db}, nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	return d.db.Close()
}

// DB 获取原始数据库连接
func (d *Database) DB() *sql.DB {
	return d.db
}

// Ping 检查数据库连接
func (d *Database) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// enablePragmas 启用SQLite优化选项
func enablePragmas(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA cache_size = 1000000000",
		"PRAGMA foreign_keys = true",
		"PRAGMA temp_store = memory",
		"PRAGMA busy_timeout = 5000",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute pragma %s: %w", pragma, err)
		}
	}

	return nil
}

// Transaction 事务管理器
type Transaction struct {
	tx *sql.Tx
}

// BeginTransaction 开始事务
func (d *Database) BeginTransaction(ctx context.Context) (*Transaction, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Transaction{tx: tx}, nil
}

// Commit 提交事务
func (t *Transaction) Commit() error {
	return t.tx.Commit()
}

// Rollback 回滚事务
func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}

// Exec 执行SQL语句
func (t *Transaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

// Query 查询
func (t *Transaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.Query(query, args...)
}

// QueryRow 查询单行
func (t *Transaction) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRow(query, args...)
}
