package resp

import (
	"fmt"
	"net"
)

type RedisCmd struct {
	Command string
	Args    []string
}

func DecodeArrayString(data []byte) ([]string, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}

	ts := value.([]any)
	token := make([]string, len(ts))
	for i := range token {
		token[i] = ts[i].(string)
	}
	return token, nil
}

func EvalAndRespond(cmd *RedisCmd, conn net.Conn) error {
	switch cmd.Command {
	case "PING":
		return EvalPing(cmd.Args, conn)
	default:
		return EvalPing(cmd.Args, conn)
	}
}

func EvalPing(args []string, conn net.Conn) error {
	var b []byte

	if len(args) >= 2 {
		conn.Write([]byte("-ERR invalid number of arguments for PING command\r\n"))
	}

	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}

	_, err := conn.Write(b)
	return err
}

func Encode(value any, isSimple bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimple {
			return []byte(fmt.Sprintf("+%s\r\n", v))
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
	}
	return []byte{}
}
