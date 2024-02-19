package main

import (
	"bytes"
	"io"
	"os"
	"syscall/js"
)

// https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
func captureStdOut() string {
    old := os.Stdout // keep backup of the real stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    run()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	return <-outC


}

func wrapperFunc() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) string {
		return captureStdOut()
	})
}

func main() {
	ch := make(chan struct{}, 0)
	js.Global().Set("wrapperFunc", wrapperFunc())
	<-ch
}
