package migration

import "github.com/bestruirui/bestsub/internal/database/migration"

// Migration001Table 初始数据库架构
func Migration001Table() string {
	return `
CREATE TABLE IF NOT EXISTS "save_template_relations" (
	"save_id" INTEGER NOT NULL,
	"template_id" INTEGER NOT NULL,
	PRIMARY KEY("save_id", "template_id"),
	FOREIGN KEY ("save_id") REFERENCES "sub_save"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("template_id") REFERENCES "sub_output_templates"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "auth" (
	"id" INTEGER,
	"user_name" TEXT NOT NULL UNIQUE,
	"password" TEXT NOT NULL,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "system_config" (
	"key" TEXT NOT NULL UNIQUE,
	"group_name" TEXT NOT NULL,
	"value" TEXT,
	"description" TEXT,
	PRIMARY KEY("key")
);

CREATE TABLE IF NOT EXISTS "notify_templates" (
	"id" INTEGER,
	"name" TEXT,
	"description" TEXT,
	"templates" TEXT NOT NULL,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "notify_config" (
	"id" INTEGER NOT NULL UNIQUE,
	"enable" BOOLEAN NOT NULL,
	"name" TEXT,
	"description" TEXT,
	"type" TEXT NOT NULL,
	"config" TEXT NOT NULL,
	"test_result" TEXT,
	"last_test" DATETIME,
	"created_at" DATETIME NOT NULL,
	"updated_at" DATETIME NOT NULL,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "tasks" (
	"id" INTEGER,
	"enable" BOOLEAN NOT NULL,
	"name" TEXT,
	"description" TEXT,
	"is_sys_task" BOOLEAN,
	"cron" TEXT,
	"timeout" INTEGER NOT NULL,
	"type" TEXT NOT NULL,
	"log_level" TEXT,
	"config" TEXT,
	"retry" INTEGER,
	"last_run_result" TEXT,
	"last_run_time" DATETIME,
	"last_run_duration" INTEGER,
	"success_count" INTEGER,
	"failed_count" INTEGER,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "storage_configs" (
	"id" INTEGER,
	"enable" BOOLEAN NOT NULL DEFAULT true,
	"name" TEXT,
	"description" TEXT,
	"type" TEXT NOT NULL,
	"config" TEXT NOT NULL,
	"test_result" TEXT,
	"last_test" DATETIME,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "sub_output_templates" (
	"id" INTEGER,
	"enable" BOOLEAN NOT NULL DEFAULT true,
	"name" TEXT,
	"description" TEXT,
	"type" TEXT NOT NULL,
	"template" TEXT NOT NULL,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "sub_node_filter_rules" (
	"id" INTEGER,
	"enable" BOOLEAN,
	"name" TEXT,
	"description" TEXT,
	"field" TEXT NOT NULL,
	"operator" TEXT NOT NULL,
	"value" TEXT NOT NULL,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "subs" (
	"id" INTEGER,
	"enable" BOOLEAN NOT NULL DEFAULT true,
	"name" TEXT,
	"description" TEXT,
	"url" TEXT NOT NULL,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "sub_save" (
	"id" INTEGER,
	"enable" BOOLEAN NOT NULL DEFAULT false,
	"name" TEXT,
	"description" TEXT,
	"rename" TEXT,
	"file_name" TEXT NOT NULL,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "sub_share_links" (
	"id" INTEGER,
	"enable" BOOLEAN NOT NULL DEFAULT false,
	"name" TEXT NOT NULL,
	"description" TEXT,
	"rename" TEXT,
	"access_count" INTEGER NOT NULL DEFAULT 0,
	"max_access_count" INTEGER NOT NULL,
	"token" TEXT NOT NULL UNIQUE,
	"expires" DATETIME NOT NULL,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "save_filter_relations" (
	"save_id" INTEGER NOT NULL,
	"filter_id" INTEGER NOT NULL,
	PRIMARY KEY("save_id", "filter_id"),
	FOREIGN KEY ("save_id") REFERENCES "sub_save"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("filter_id") REFERENCES "sub_node_filter_rules"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "save_storage_relations" (
	"save_id" INTEGER NOT NULL,
	"storage_id" INTEGER NOT NULL,
	PRIMARY KEY("save_id", "storage_id"),
	FOREIGN KEY ("save_id") REFERENCES "sub_save"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("storage_id") REFERENCES "storage_configs"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "share_template_relations" (
	"share_id" INTEGER NOT NULL,
	"template_id" INTEGER NOT NULL,
	PRIMARY KEY("share_id", "template_id"),
	FOREIGN KEY ("share_id") REFERENCES "sub_share_links"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("template_id") REFERENCES "sub_output_templates"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "share_filter_relations" (
	"share_id" INTEGER NOT NULL,
	"filter_id" INTEGER NOT NULL,
	PRIMARY KEY("share_id", "filter_id"),
	FOREIGN KEY ("share_id") REFERENCES "sub_share_links"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("filter_id") REFERENCES "sub_node_filter_rules"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "share_sub_relations" (
	"share_id" INTEGER NOT NULL,
	"sub_id" INTEGER NOT NULL,
	PRIMARY KEY("share_id", "sub_id"),
	FOREIGN KEY ("share_id") REFERENCES "sub_share_links"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("sub_id") REFERENCES "subs"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "save_sub_relations" (
	"save_id" INTEGER NOT NULL,
	"sub_id" INTEGER NOT NULL,
	PRIMARY KEY("save_id", "sub_id"),
	FOREIGN KEY ("sub_id") REFERENCES "subs"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("save_id") REFERENCES "sub_save"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "sub_task_relations" (
	"sub_id" INTEGER NOT NULL,
	"task_id" INTEGER NOT NULL,
	PRIMARY KEY("sub_id", "task_id"),
	FOREIGN KEY ("sub_id") REFERENCES "subs"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("task_id") REFERENCES "tasks"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "save_task_relations" (
	"save_id" INTEGER NOT NULL,
	"task_id" INTEGER NOT NULL,
	PRIMARY KEY("save_id", "task_id"),
	FOREIGN KEY ("task_id") REFERENCES "tasks"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("save_id") REFERENCES "sub_save"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "notify_task_relations" (
	"task_id" INTEGER NOT NULL,
	"notify_id" INTEGER NOT NULL,
	PRIMARY KEY("task_id", "notify_id"),
	FOREIGN KEY ("task_id") REFERENCES "tasks"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("notify_id") REFERENCES "notify_config"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "task_notify_template_relations" (
	"task_id" INTEGER NOT NULL,
	"notify_template_id" INTEGER NOT NULL,
	PRIMARY KEY("task_id", "notify_template_id"),
	FOREIGN KEY ("notify_template_id") REFERENCES "notify_templates"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("task_id") REFERENCES "tasks"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "migrations" (
	"date" INTEGER NOT NULL UNIQUE,
	"version" TEXT NOT NULL,
	"description" TEXT NOT NULL,
	"applied_at" DATETIME NOT NULL,
	PRIMARY KEY("date")
);
`
}

// init 自动注册迁移
func init() {
	migration.Register(ClientName, 202507171100, "dev", "Tables", Migration001Table)
}
