package transformer

import "strconv"

// repeatedly calling the function returned will give unique names with the given prefix
func getUniqueCounter(ident string) func() string {
	var argN = 0
	return func() string {
		str := ident + strconv.Itoa(argN)
		argN++
		return str
	}
}
