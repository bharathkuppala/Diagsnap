package utility

import "strconv"

// ConverAtoI ...
func ConverAtoI(data string) (conStr int, err error) {
	conStr, err = strconv.Atoi(data)
	if err != nil {
		return 0, err
	}
	return conStr, nil
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
