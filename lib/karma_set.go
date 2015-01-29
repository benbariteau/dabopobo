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
	return fmt.Sprintf("(%v++, %v--, %v+-)", k.plusplus, k.minusminus, k.plusminus)
}
