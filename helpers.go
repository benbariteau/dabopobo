package main

// canonicalizeSuffix returns the "canonical" version of mutation suffixes
func canonicalizeSuffix(suffix string) string {
	switch suffix {
	case "--", "++", "+-":
		return suffix
	case "-+":
		return "+-"
	default:
		return suffix
	}
}
