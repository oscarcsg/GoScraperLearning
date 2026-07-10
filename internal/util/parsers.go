package util

import "fmt"

func ParseUint8(txt string) (uint8, error) {
	if !CheckOnlyNumber(txt) {
		return 0, fmt.Errorf("not_only_num")
	}
		
    var n uint8
    for i := 0; i < len(txt); i++ {
        n = n*10 + (txt[i] - '0')
    }
    return n, nil
}

func ParseUint16(txt string) (uint16, error) {
	if !CheckOnlyNumber(txt) {
		return 0, fmt.Errorf("not_only_num")
	}

	var n uint16
	for i := 0; i < len(txt); i++ {
		n = n*10 + uint16(txt[i] - '0')
	}
	return n, nil
}
