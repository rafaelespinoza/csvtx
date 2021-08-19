package entity

import "fmt"

type AmountSubunits int64

func (a AmountSubunits) String() string {
	absVal := a
	if a < 0 {
		absVal *= -1
	}
	dollars := absVal / 100
	cents := absVal % 100
	return fmt.Sprintf("%d.%02d", dollars, cents)
}
