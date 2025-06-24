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
	pragmas := map[string]string{
		"journal_mode":       "WAL",    // 启用WAL模式提高并发性能
		"synchronous":        "NORMAL", // 平衡性能和安全性
		"cache_size":         "-64000", // 64MB缓存
		"foreign_keys":       "ON",     // 启用外键约束
		"temp_store":         "MEMORY", // 临时表存储在内存中
		"busy_timeout":       "5000",   // 5秒忙等待超时
		"wal_autocheckpoint": "1000",   // WAL自动检查点
		"optimize":           "",       // 优化数据库
	}

	for key, value := range pragmas {
		var query string
		if value == "" {
			query = fmt.Sprintf("PRAGMA %s", key)
		} else {
			query = fmt.Sprintf("PRAGMA %s = %s", key, value)
		}

		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute pragma %s: %w", query, err)
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

// Stats 获取数据库统计信息
func (d *Database) Stats() sql.DBStats {
	return d.db.Stats()
}

// Exec 执行SQL语句
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

// ExecContext 执行SQL语句（带上下文）
func (d *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

// Query 查询
func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

// QueryContext 查询（带上下文）
func (d *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

// QueryRow 查询单行
func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(query, args...)
}

// QueryRowContext 查询单行（带上下文）
func (d *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

// Vacuum 清理数据库
func (d *Database) Vacuum(ctx context.Context) error {
	_, err := d.db.ExecContext(ctx, "VACUUM")
	return err
}

// Analyze 分析数据库统计信息
func (d *Database) Analyze(ctx context.Context) error {
	_, err := d.db.ExecContext(ctx, "ANALYZE")
	return err
}
