package lib

import (
	"fmt"
)

type karmaSet struct {
	plusplus   int
	minusminus int
	plusminus  int
}

func (k karmaSet) value() int {
	return k.plusplus - k.minusminus
}

func (k karmaSet) String() string {
	ratio := float64(k.plusplus) / float64(k.minusminus)
	return fmt.Sprintf("(%v++, %v--, %v+-) +/- ratio: %.3g", k.plusplus, k.minusminus, k.plusminus, ratio)
}
