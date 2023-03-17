package ultilities

import "time"

func TimeInUTC(t time.Time) time.Time {
	return t.UTC()
}
func TimeInLocal(t time.Time) time.Time {
	return t.In(time.Local)
}
