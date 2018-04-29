package csvtx

import (
	"strconv"
)

type Amount int

func (a Amount) String() string {
	d := int(a / 100) // dollars
	c := int(a % 100) // cents
	s := numStr(d) + "." + cents(c)

	if a.isNegative() {
		return "(" + s + ")"
	} else {
		return s
	}
}

func (a *Amount) isNegative() bool {
	return int(*a) < 0
}

func cents(c int) string {
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
