package main

import (
	"fmt"
	"strings"
)

// canonicalizeSuffix returns the "canonical" version of mutation suffixes
func canonicalizeSuffix(suffix string) string {
	if len(suffix) > 0 {
		suffix = suffix[:2]
	}
	switch suffix {
	case "--", "++", "+-":
		return suffix
	case "-+":
		return "+-"
	default:
		return suffix
	}
}

func maybeRemoveParens(s string) string {
	if s == "" {
		return s
	} else if s[0] == '(' && s[len(s)-1] == ')' {
		return s[1 : len(s)-1]
	}
	return s
}

func filterMutations(mutations [][]string, filters ...string) (mutationList []karmaMutation) {
	mutationSet := make(map[karmaMutation]bool)

	filterSet := make(map[string]bool)
	for _, filter := range filters {
		filterSet[filter] = true
	}

	for _, mutation := range mutations {
		m := newKarmaMutation(mutation[1], mutation[2])

		if m.identifier != "" && !filterSet[m.identifier] && !mutationSet[m] {
			mutationList = append(mutationList, m)
			mutationSet[m] = true
		}
	}
	return
}

func newKarmaMutation(identifier, op string) karmaMutation {
	return karmaMutation{
		identifier: strings.ToLower(maybeRemoveParens(identifier)),
		op:         canonicalizeSuffix(op),
	}
}

type karmaMutation struct {
	identifier string
	op         string
}

func (m karmaMutation) String() string {
	return m.key()
}

func (m karmaMutation) key() string {
	return fmt.Sprintf("%v%v", m.identifier, m.op)
}
