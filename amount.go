package csvtx

import (
	"strconv"
)

type Amount int

func (a Amount) String() string {
	absVal := a.abs()
	dollars := absVal / 100
	cents := absVal % 100
	return numStr(dollars) + "." + centStr(cents)
}

func (a *Amount) isNegative() bool {
	return int(*a) < 0
}

func (a Amount) abs() int {
	num := int(a)

	if a.isNegative() {
		return num * -1
	} else {
		return num
	}
}

func centStr(c int) string {
	s := numStr(c)

	if c < 10 {
		return "0" + s
	} else {
		return s
	}
}

func numStr(n int) string {
	return strconv.FormatInt(int64(n), 10)
}
