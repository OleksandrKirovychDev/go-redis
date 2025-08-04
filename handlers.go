package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET": set,
	"GET": get,
	"HSET": hSet,
	"HGET": hGet,
	"HGETALL": hGetAll,
}

func ping(args []Value) Value {
	return Value{typ: "string", str: "PONG"}
}


var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if(len(args) != 2) {
		return Value{typ: "error", str: "wrong number of arguments for SET command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if(len(args) != 1) {
		return Value{typ: "error", str: "wrong number of arguments for GET command"}
	}

	key := args[0].bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hSet(args []Value) Value {
	if(len(args) != 3) {
		return Value{typ: "error", str: "wrong number of arguments for HSET command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hGet(args []Value) Value {
	if(len(args) != 2) {
		return Value{typ: "error", str: "wrong number of arguments for HGET command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hGetAll(args []Value) Value {
	if(len(args) != 1) {
		return Value{typ: "error", str: "wrong number of arguments for HGET command"}
	}

	hash := args[0].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "array", array: []Value{}}
	}

	result := make([]Value, 0, len(value))

	for k, v := range value {
		result = append(result, Value{typ: "bulk", bulk: k})
		result = append(result, Value{typ: "bulk", bulk: v})
	}

	return Value{typ: "array", array: result}
}
