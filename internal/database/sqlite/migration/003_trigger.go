package migration

import "github.com/bestruirui/bestsub/internal/database/migration"

func Migration003Trigger() string {
	return `
CREATE TRIGGER IF NOT EXISTS delete_sub_link_tasks
AFTER DELETE ON sub_links
FOR EACH ROW
BEGIN
    DELETE FROM tasks 
    WHERE id IN (
        SELECT task_id FROM sub_task_relations 
        WHERE sub_id = OLD.id
    );
END;

CREATE TRIGGER IF NOT EXISTS delete_save_config_tasks
AFTER DELETE ON sub_save_configs
FOR EACH ROW
BEGIN
    DELETE FROM tasks
    WHERE id IN (
        SELECT task_id FROM save_task_relations
        WHERE save_config_id = OLD.id
    );
END;
	`
}

func init() {
	migration.Register(migrations, "003", "Triggers", Migration003Trigger)
}
