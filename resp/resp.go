package resp

import "fmt"

func Decode(data []byte) (any, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("no data")
	}
	val, _, err := DecodeOne(data)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func DecodeOne(data []byte) (any, int, error) {
	if len(data) == 0 {
		return 0, 0, fmt.Errorf("no data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '$':
		return readBulkString(data)
	case '-':
		return readError(data)
	case ':':
		return readInteger(data)
	case '*':
		return readArray(data)

	}
	return 0, 0, nil
}

func readLength(data []byte) (int, int) {
	length := 0

	for i, b := range data {
		if b < '0' || b > '9' {
			return length, i + 2
		}
		length = length*10 + int(b-'0')
	}

	return 0, 0
}

func readArray(data []byte) (any, int, error) {
	pos := 1

	count, delta := readLength(data[pos:])
	pos += delta

	var elements []interface{} = make([]interface{}, count)
	for i := range elements {
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elements[i] = elem
		pos += delta
	}
	return elements, pos, nil
}

func readInteger(data []byte) (int64, int, error) {
	pos := 1
	var value int64 = 0
	for ; data[pos] != 'r'; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}
	return value, pos + 2, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

func readBulkString(data []byte) (string, int, error) {
	pos := 1

	len, delta := readLength(data[pos:])
	pos += delta

	return string(data[pos:(pos + len)]), pos + len + 2, nil
}

func readSimpleString(data []byte) (string, int, error) {
	if len(data) == 0 {
		return "", 0, fmt.Errorf("no data")
	}
	pos := 1

	for ; data[pos] != '\r'; pos++ {
	}

	return string(data[1:pos]), pos + 2, nil
}
