package tools

import "strconv"

func StringToUint(str string) (uint, error) {
	temporary, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(temporary), nil
}