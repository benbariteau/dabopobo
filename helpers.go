package main

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
