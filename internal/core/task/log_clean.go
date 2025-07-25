package task

import (
	"time"

	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

func StartLogClean() {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := log.CleanupOldLogs(op.GetConfigInt("log.retention_days")); err != nil {
				log.Errorf("Scheduled log cleanup failed: %v", err)
			}
		}
	}()
}
