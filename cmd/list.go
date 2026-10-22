package cmd

import (
	"redigo/resp"
	"redigo/store"
	"strconv"
)

func lpush(env resp.Envelope, h *Handler) string {
	if env.Size < 3 {
		return HandleErr(WNA, env)
	}

	key := env.Array[1].String
	values := make([]string, env.Size-2)

	for i, e := range env.Array {
		if i > 1 {
			values[i-2] = e.String
		}
	}

	size, err := h.store.LPUSH(key, values)

	if err != nil {
		if err == store.ErrWrongType {
			return HandleErr(WTO, env)
		}
	}

	resEnv := resp.Envelope{
		OpCode:  resp.Integer,
		Integer: size,
	}

	return resp.FormatMapper(resEnv)

}

func rpush(env resp.Envelope, h *Handler) string {
	if env.Size < 3 {
		return HandleErr(WNA, env)
	}

	key := env.Array[1].String
	values := make([]string, env.Size-2)

	for i, e := range env.Array {
		if i > 1 {
			values[i-2] = e.String
		}
	}

	size, err := h.store.RPUSH(key, values)

	if err != nil {
		if err == store.ErrWrongType {
			return HandleErr(WTO, env)
		}
	}

	resEnv := resp.Envelope{
		OpCode:  resp.Integer,
		Integer: size,
	}

	return resp.FormatMapper(resEnv)

}

func lpop(env resp.Envelope, h *Handler) string {
	if env.Size < 2 || env.Size > 3 {
		return HandleErr(WNA, env)
	}

	key := env.Array[1].String

	var count int = 1
	var err error

	if env.Size == 3 {
		count, err = strconv.Atoi(env.Array[2].String)
		if err != nil {
			return HandleErr(NAI, env)
		}
	}

	values, err := h.store.LPOP(key, count)

	if err != nil {
		if err == store.ErrWrongType {
			return HandleErr(WTO, env)
		}
	}

	if len(values) == 0 {
		return NIL()
	}

	if len(values) == 1 {
		resEnv := resp.Envelope{
			OpCode: resp.BulkString,
			String: values[0],
			Size:   len(values[0]),
		}

		return resp.FormatMapper(resEnv)
	}

	resEnv := resp.Envelope{
		OpCode: resp.Array,
		Size:   len(values),
		Array:  make([]resp.Envelope, len(values)),
	}

	for i, v := range values {
		valEnv := resp.Envelope{
			OpCode: resp.BulkString,
			String: v,
			Size:   len(v),
		}

		resEnv.Array[i] = valEnv

	}

	return resp.FormatMapper(resEnv)

}

func rpop(env resp.Envelope, h *Handler) string {
	if env.Size < 2 || env.Size > 3 {
		return HandleErr(WNA, env)
	}

	key := env.Array[1].String

	var count int = 1
	var err error

	if env.Size == 3 {
		count, err = strconv.Atoi(env.Array[2].String)
		if err != nil {
			return HandleErr(NAI, env)
		}
	}

	values, err := h.store.RPOP(key, count)

	if err != nil {
		if err == store.ErrWrongType {
			return HandleErr(WTO, env)
		}
	}

	if len(values) == 0 {
		return NIL()
	}

	if len(values) == 1 {
		resEnv := resp.Envelope{
			OpCode: resp.BulkString,
			String: values[0],
			Size:   len(values[0]),
		}

		return resp.FormatMapper(resEnv)
	}

	resEnv := resp.Envelope{
		OpCode: resp.Array,
		Size:   len(values),
		Array:  make([]resp.Envelope, len(values)),
	}

	for i, v := range values {
		valEnv := resp.Envelope{
			OpCode: resp.BulkString,
			String: v,
			Size:   len(v),
		}

		resEnv.Array[i] = valEnv

	}

	return resp.FormatMapper(resEnv)

}

func lrange(env resp.Envelope, h *Handler) string {
	if env.Size != 4 {
		return HandleErr(WNA, env)
	}

	key := env.Array[1].String
	start, err := strconv.Atoi(env.Array[2].String)
	if err != nil {
		return HandleErr(store.ErrValNotInt, env)
	}
	stop, err := strconv.Atoi(env.Array[3].String)
	if err != nil {
		return HandleErr(store.ErrValNotInt, env)
	}

	values, err := h.store.LRANGE(key, start, stop)
	if err != nil {
		if err == store.ErrWrongType {
			return HandleErr(WTO, env)
		}
	}

	valEnv := make([]resp.Envelope, len(values))
	for i, val := range values {
		valEnv[i] = resp.Envelope{
			OpCode: resp.BulkString,
			Size:   len(val),
			String: val,
		}
	}

	resEnv := resp.Envelope{
		OpCode: resp.Array,
		Size:   len(values),
		Array:  valEnv,
	}

	return resp.FormatMapper(resEnv)
}

func llen(env resp.Envelope, h *Handler) string {
	if env.Size != 2 {
		return HandleErr(WNA, env)
	}

	key := env.Array[1].String
	size, err := h.store.LLEN(key)

	if err != nil {
		if err == store.ErrWrongType {
			return HandleErr(WTO, env)
		}
	}

	resEnv := resp.Envelope{
		OpCode:  resp.Integer,
		Integer: size,
	}

	return resp.FormatMapper(resEnv)
}
