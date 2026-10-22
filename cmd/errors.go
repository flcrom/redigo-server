package cmd

import (
	"errors"
	"redigo/resp"
)

var UNC error = errors.New("UNC") // unknown command with arf
var WNA error = errors.New("WNA") // wrong number of arguments
var NAI error = errors.New("NAI") // not an integer
var WTO error = errors.New("WTO") // wrong type operation
var NCP error = errors.New("NCP") // no command present

type ErrData struct {
	errCode string
	errCmd  string
	args    []string
}

func HandleErr(errCode error, errEnv resp.Envelope) string {

	var errType resp.Type = resp.SimpleError
	errStr := "Default Err Str"

	switch errCode {
	case UNC:
		var args string
		for i := range errEnv.Size {
			args = args + "'" + errEnv.Array[i].String + "' "
		}
		errStr = "ERR unknown command '" + errEnv.Array[0].String + "', with args beginning with: " + args

	case WNA:
		errStr = "ERR wrong number of arguments for '" + errEnv.Array[0].String + "' command"

	case NAI:
		errStr = "ERR value is not an integer or out of range"

	case WTO:
		errStr = "WRONGTYPE Operation against a key holding the wrong kind of value"

	}

	env := resp.Envelope{
		OpCode: errType,
		String: errStr,
		Size:   len(errStr),
	}

	return resp.FormatMapper(env)
}
