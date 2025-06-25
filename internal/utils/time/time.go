package time

import timestd "time"

func Now() timestd.Time {
	return timestd.Now().Local()
}
