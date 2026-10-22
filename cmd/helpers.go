package cmd

import (
	"redigo/resp"
)

func NIL() string {
	NIL := resp.Envelope{
		OpCode: resp.BulkString,
		Size:   -1,
	}

	return resp.FormatMapper(NIL)
}
