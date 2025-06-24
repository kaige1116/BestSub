package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bestruirui/bestsub/internal/database/sqlite"
)

func main() {
	// 创建临时数据库文件
	dbPath := "/tmp/test_migration.db"
	
	// 删除已存在的数据库文件
	os.Remove(dbPath)
	
	// 创建数据库连接
	db, err := sqlite.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()
	
	// 创建仓库
	repo := sqlite.NewRepository(db)
	
	ctx := context.Background()
	
	// 测试迁移功能
	fmt.Println("开始测试迁移功能...")
	
	// 运行迁移
	fmt.Println("运行迁移...")
	if err := repo.RunMigrations(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	fmt.Println("迁移完成!")
	
	// 获取迁移状态
	fmt.Println("获取迁移状态...")
	status, err := repo.GetMigrationStatus(ctx)
	if err != nil {
		log.Fatalf("Failed to get migration status: %v", err)
	}
	
	fmt.Printf("迁移状态:\n")
	fmt.Printf("  总迁移数: %d\n", status.TotalMigrations)
	fmt.Printf("  已应用迁移数: %d\n", status.AppliedMigrations)
	fmt.Printf("  待执行迁移数: %d\n", status.PendingMigrations)
	
	if status.LastMigration != nil {
		fmt.Printf("  最后迁移: %s (%s)\n", status.LastMigration.ID, status.LastMigration.Description)
	}
	
	fmt.Println("\n已应用的迁移:")
	for _, migration := range status.AppliedList {
		fmt.Printf("  - %s: %s (应用时间: %s)\n", migration.ID, migration.Description, migration.AppliedAt.Format("2006-01-02 15:04:05"))
	}
	
	fmt.Println("\n待执行的迁移:")
	for _, migration := range status.PendingList {
		fmt.Printf("  - %s: %s\n", migration.ID(), migration.Description())
	}
	
	// 测试数据库表是否创建成功
	fmt.Println("\n测试数据库表...")
	tables := []string{"auth", "sessions", "system_config", "tasks", "sub_links"}
	for _, table := range tables {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='%s'", table)
		err := db.QueryRowContext(ctx, query).Scan(&count)
		if err != nil {
			log.Printf("Failed to check table %s: %v", table, err)
			continue
		}
		if count > 0 {
			fmt.Printf("  ✓ 表 %s 存在\n", table)
		} else {
			fmt.Printf("  ✗ 表 %s 不存在\n", table)
		}
	}
	
	// 测试索引是否创建成功
	fmt.Println("\n测试数据库索引...")
	indexes := []string{"idx_sessions_token_hash", "idx_tasks_status", "idx_sub_links_enabled"}
	for _, index := range indexes {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='%s'", index)
		err := db.QueryRowContext(ctx, query).Scan(&count)
		if err != nil {
			log.Printf("Failed to check index %s: %v", index, err)
			continue
		}
		if count > 0 {
			fmt.Printf("  ✓ 索引 %s 存在\n", index)
		} else {
			fmt.Printf("  ✗ 索引 %s 不存在\n", index)
		}
	}
	
	fmt.Println("\n迁移测试完成!")
	
	// 清理测试文件
	os.Remove(dbPath)
}
