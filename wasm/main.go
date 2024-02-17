package main

import (
	"syscall/js"
)

func wrapperFunc() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		run()
		return nil
	})
}

func main() {
	ch := make(chan struct{}, 0)
	js.Global().Set("wrapperFunc", wrapperFunc())
	<-ch
}
