package main

import (
	"github.com/yuin/gopher-lua"
	"reflect"
	"fmt"
	"math/rand"
	"time"
)

type (
	Variables struct {
		N int `lua:"n"`
	}

	RuntimeRegistry struct {
		Variables `lua:"variables"`
	}
)

var (
	luaScript = `
		function add ()
			return runtime.variables.n + 100;
		end
	`

	runtime = RuntimeRegistry{
		Variables{
			N: rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100),
		},
	}
)

func main() {
	L := lua.NewState()
	defer L.Close()

	// creating Lua table from Go structure
	L.SetGlobal(
		"runtime",
		populateLuaTableWith(runtime, L, L.NewTable()),
	)

	// executing Lua script
	if err := L.DoString(luaScript); err != nil {
		panic(err)
	}

	// getting result of `logical_expression_result` variable in Go
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("add"),
		NRet:    1,
		Protect: true,
	}); err != nil {
		panic(err)
	}

	if result, ok := L.Get(-1).(lua.LNumber); ok {
		fmt.Println("add() result:", result)
	}
}

// populateLuaTableWith take an instance of any interface, inspects it
// and populates given Lua table with values from this interface
func populateLuaTableWith(structure interface{}, state *lua.LState, table lua.LValue) lua.LValue {
	var value = reflect.ValueOf(structure)
	for i := 0; i < value.NumField(); i++ {
		var (
			field = value.Field(i)
			tag   = value.Type().Field(i).Tag.Get("lua")
			value lua.LValue
		)

		switch field.Kind() {
		case reflect.Struct:
			value = populateLuaTableWith(field.Interface(), state, table)
		case reflect.Int:
			value = lua.LNumber(field.Int())
		}
		state.SetField(table, tag, value)
	}

	return table
}