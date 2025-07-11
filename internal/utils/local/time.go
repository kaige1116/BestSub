package local

import time "time"

func Time() time.Time {
	return time.Now().Local()
}
