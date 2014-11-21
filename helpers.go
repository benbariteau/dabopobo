package main

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
