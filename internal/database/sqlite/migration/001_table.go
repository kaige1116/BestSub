package migration

import "github.com/bestruirui/bestsub/internal/database/migration"

// Migration001Table 初始数据库架构
func Migration001Table() string {
	return `
CREATE TABLE IF NOT EXISTS "auth" (
	"id" INTEGER,
	"user_name" TEXT NOT NULL UNIQUE,
	"password" TEXT NOT NULL,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "sessions" (
	"id" INTEGER,
	"is_active" BOOLEAN NOT NULL DEFAULT true,
	"ip_address" TEXT,
	"user_agent" TEXT,
	"expires_at" DATETIME NOT NULL,
	"token_hash" TEXT NOT NULL UNIQUE,
	"refresh_token" TEXT UNIQUE,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "system_config" (
	"id" INTEGER,
	"group_name" TEXT NOT NULL,
	"type" TEXT NOT NULL,
	"key" TEXT NOT NULL,
	"value" TEXT,
	"description" TEXT,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "notify" (
	"id" INTEGER,
	"enable" BOOLEAN NOT NULL DEFAULT true,
	"name" TEXT,
	"type" TEXT NOT NULL,
	"config" TEXT NOT NULL,
	"test_result" TEXT,
	"last_test" DATETIME,
	"description" TEXT,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "tasks" (
	"id" INTEGER,
	"enable" BOOLEAN NOT NULL,
	"name" TEXT,
	"cron" TEXT,
	"type" TEXT NOT NULL,
	"status" TEXT NOT NULL DEFAULT 'pending',
	"config" TEXT,
	"last_run_result" TEXT,
	"last_run_time" DATETIME,
	"last_run_duration" INTEGER,
	"description" TEXT,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "sub_storage_configs" (
	"id" INTEGER,
	"enable" BOOLEAN NOT NULL DEFAULT true,
	"name" TEXT,
	"type" TEXT NOT NULL,
	"config" TEXT NOT NULL,
	"test_result" TEXT,
	"last_test" DATETIME,
	"description" TEXT,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "sub_output_templates" (
	"id" INTEGER,
	"enable" BOOLEAN NOT NULL DEFAULT true,
	"name" TEXT,
	"type" TEXT NOT NULL,
	"template" TEXT NOT NULL,
	"description" TEXT,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "sub_node_filter_rules" (
	"id" INTEGER,
	"name" TEXT,
	"field" TEXT NOT NULL,
	"operator" TEXT NOT NULL,
	"value" TEXT NOT NULL,
	"description" TEXT,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "sub_links" (
	"id" INTEGER,
	"enable" BOOLEAN NOT NULL DEFAULT true,
	"name" TEXT,
	"url" TEXT NOT NULL,
	"description" TEXT,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "sub_save_configs" (
	"id" INTEGER,
	"enable" BOOLEAN NOT NULL DEFAULT false,
	"name" TEXT,
	"rename" TEXT,
	"type" TEXT NOT NULL,
	"file_name" TEXT NOT NULL,
	"description" TEXT,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "sub_share_links" (
	"id" INTEGER,
	"enable" BOOLEAN NOT NULL DEFAULT false,
	"name" TEXT,
	"rename" TEXT,
	"access_count" INTEGER NOT NULL DEFAULT 0,
	"max_access_count" INTEGER NOT NULL,
	"token" TEXT NOT NULL UNIQUE,
	"expires" DATETIME NOT NULL,
	"description" TEXT,
	"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "save_template_relations" (
	"save_config_id" INTEGER NOT NULL,
	"template_id" INTEGER NOT NULL,
	PRIMARY KEY("save_config_id", "template_id"),
	FOREIGN KEY ("save_config_id") REFERENCES "sub_save_configs"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("template_id") REFERENCES "sub_output_templates"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "save_fitter_relations" (
	"save_config_id" INTEGER NOT NULL,
	"fitter_id" INTEGER NOT NULL,
	PRIMARY KEY("save_config_id", "fitter_id"),
	FOREIGN KEY ("save_config_id") REFERENCES "sub_save_configs"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("fitter_id") REFERENCES "sub_node_filter_rules"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "save_storage_relations" (
	"save_config_id" INTEGER NOT NULL,
	"storage_id" INTEGER NOT NULL,
	PRIMARY KEY("save_config_id", "storage_id"),
	FOREIGN KEY ("save_config_id") REFERENCES "sub_save_configs"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("storage_id") REFERENCES "sub_storage_configs"("id")
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

CREATE TABLE IF NOT EXISTS "share_fitter_relations" (
	"share_id" INTEGER NOT NULL,
	"fitter_id" INTEGER NOT NULL,
	PRIMARY KEY("share_id", "fitter_id"),
	FOREIGN KEY ("share_id") REFERENCES "sub_share_links"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("fitter_id") REFERENCES "sub_node_filter_rules"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "share_sub_relations" (
	"share_id" INTEGER NOT NULL,
	"sub_id" INTEGER NOT NULL,
	PRIMARY KEY("share_id", "sub_id"),
	FOREIGN KEY ("share_id") REFERENCES "sub_share_links"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("sub_id") REFERENCES "sub_links"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "save_sub_relations" (
	"save_config_id" INTEGER NOT NULL,
	"sub_id" INTEGER NOT NULL,
	PRIMARY KEY("save_config_id", "sub_id"),
	FOREIGN KEY ("sub_id") REFERENCES "sub_links"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("save_config_id") REFERENCES "sub_save_configs"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "sub_task_relations" (
	"sub_id" INTEGER NOT NULL,
	"task_id" INTEGER NOT NULL,
	PRIMARY KEY("sub_id", "task_id"),
	FOREIGN KEY ("sub_id") REFERENCES "sub_links"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("task_id") REFERENCES "tasks"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "save_task_relations" (
	"save_config_id" INTEGER NOT NULL,
	"task_id" INTEGER NOT NULL,
	PRIMARY KEY("save_config_id", "task_id"),
	FOREIGN KEY ("task_id") REFERENCES "tasks"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("save_config_id") REFERENCES "sub_save_configs"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "notify_task_relations" (
	"task_id" INTEGER NOT NULL,
	"notify_id" INTEGER NOT NULL,
	PRIMARY KEY("task_id", "notify_id"),
	FOREIGN KEY ("task_id") REFERENCES "tasks"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE,
	FOREIGN KEY ("notify_id") REFERENCES "notify"("id")
	ON UPDATE NO ACTION ON DELETE CASCADE
);
`
}

// init 自动注册所有迁移
func init() {
	migration.Register(migrations, "001", "Tables", Migration001Table)
}
