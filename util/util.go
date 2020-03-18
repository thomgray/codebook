package util

// Check ...
func Check(err error) {
	if err != nil {
		panic(err)
	}
}
