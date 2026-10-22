package cmd

import (
	"bufio"
	"redigo/resp"
	"redigo/store"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	store *store.Store
}

func NewHandler(store *store.Store) *Handler {
	return &Handler{
		store: store,
	}
}

func HandleRequest(h *Handler, reader *bufio.Reader) (string, error) {
	var res string
	var env resp.Envelope

	env, err := resp.ParseInput(reader)

	if err != nil {
		return res, err
	}

	if env.Size == 0 {
		return "", NCP
	}

	res = commandRouter(env, h)
	return res, nil

}

func commandRouter(command resp.Envelope, h *Handler) string {

	op := command.Array[0]

	switch com := strings.ToUpper(op.String); com {
	case "PING":
		return ping()
	case "ECHO":
		return echo(command)
	case "SET":
		return set(command, h)
	case "GET":
		return get(command, h)
	case "INCR":
		return incr(command, h)
	case "DECR":
		return decr(command, h)
	case "LPUSH":
		return lpush(command, h)
	case "RPUSH":
		return rpush(command, h)
	case "LPOP":
		return lpop(command, h)
	case "RPOP":
		return rpop(command, h)
	case "LRANGE":
		return lrange(command, h)
	case "LLEN":
		return llen(command, h)
	default:
		return HandleErr(UNC, command)
	}

}

func ping() string {
	env := resp.Envelope{
		OpCode: resp.SimpleString,
		String: "PONG",
	}

	return resp.FormatMapper(env)
}

func echo(env resp.Envelope) string {
	if env.Size != 2 {
		return HandleErr(WNA, env)
	}

	resEnv := resp.Envelope{
		OpCode: resp.BulkString,
		Size:   env.Array[1].Size,
		String: env.Array[1].String,
	}

	return resp.FormatMapper(resEnv)
}

func set(env resp.Envelope, h *Handler) string {
	if env.Size < 3 || env.Size > 5 {
		return HandleErr(WNA, env)
	}

	key := env.Array[1].String
	value := env.Array[2].String
	var expiry time.Time

	switch env.Size {
	case 3:
		expiry = time.Time{}
	case 5:
		switch option := strings.ToUpper(env.Array[3].String); option {
		case "EX", "PX":
			interval, err := strconv.Atoi(env.Array[4].String)
			if err != nil {
				return HandleErr(NAI, env)
			}
			expiry = time.Now()

			if option == "EX" {
				expiry = expiry.Add(time.Duration(interval) * time.Second)
			} else {
				expiry = expiry.Add(time.Duration(interval) * time.Millisecond)
			}
		}
	default:
		return HandleErr(WNA, env)
	}

	h.store.SET(key, value, expiry)
	resEnv := resp.Envelope{
		OpCode: resp.SimpleString,
		String: "OK",
		Size:   2,
	}

	return resp.FormatMapper(resEnv)

}

func get(env resp.Envelope, h *Handler) string {
	if env.Size != 2 {
		return HandleErr(WNA, env)
	}
	key := env.Array[1].String
	value, ok := h.store.GET(key)

	if !ok {
		return NIL()
	}

	resEnv := resp.Envelope{
		OpCode: resp.BulkString,
		Size:   len(value),
		String: value,
	}

	return resp.FormatMapper(resEnv)
}

func incr(env resp.Envelope, h *Handler) string {
	if env.Size != 2 {
		return HandleErr(WNA, env)
	}

	key := env.Array[1].String
	value, ok := h.store.INCR(key)
	if !ok {
		return HandleErr(NAI, env)
	}

	resEnv := resp.Envelope{
		OpCode:  resp.Integer,
		Integer: value,
	}

	return resp.FormatMapper(resEnv)

}

func decr(env resp.Envelope, h *Handler) string {
	if env.Size != 2 {
		return HandleErr(WNA, env)
	}

	key := env.Array[1].String
	value, ok := h.store.DECR(key)
	if !ok {
		return HandleErr(NAI, env)
	}

	resEnv := resp.Envelope{
		OpCode:  resp.Integer,
		Integer: value,
	}

	return resp.FormatMapper(resEnv)
}
