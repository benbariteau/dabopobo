package main

import (
	"reflect"
	"testing"
)

func TestFilterMutations(t *testing.T) {
	type params struct {
		mutations [][]string
		filters   []string
	}

	testcases := []struct {
		in  params
		out []karmaMutation
	}{
		{
			params{
				[][]string{
					[]string{"butt++", "butt", "++"},
					[]string{"butt++", "butt", "++"},
				},
				[]string{"breetz"},
			},
			[]karmaMutation{
				karmaMutation{"butt", "++"},
			},
		},
		{
			params{
				[][]string{
					[]string{"butt++", "butt", "++"},
					[]string{"fart++", "fart", "++"},
				},
				[]string{"breetz"},
			},
			[]karmaMutation{
				karmaMutation{"butt", "++"},
				karmaMutation{"fart", "++"},
			},
		},
		{
			params{
				[][]string{
					[]string{"butt++", "butt", "++"},
					[]string{"breetz++", "breetz", "++"},
				},
				[]string{"breetz"},
			},
			[]karmaMutation{
				karmaMutation{"butt", "++"},
			},
		},
	}

	for _, testcase := range testcases {
		o := filterMutations(testcase.in.mutations, testcase.in.filters...)
		if !reflect.DeepEqual(o, testcase.out) {
			t.Error("Got", o, "exepected", testcase.out)
		}
	}
}
