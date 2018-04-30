package csvtx

import (
	"time"
)

// It feels overkill to create a whole package to namespace any test helper
// code. Instead, prefix any test helper methods with `th_`
func th_makeDate() time.Time {
	return time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)
}
