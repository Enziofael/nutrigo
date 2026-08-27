//go:build windows

package os

import (
	"log"
	"os"

	"golang.org/x/sys/windows"
)

func init() {
	handle := windows.Handle(os.Stdout.Fd())
	var mode uint32
	if err := windows.GetConsoleMode(handle, &mode); err == nil {
		log.Println("Enabling virtual terminal processing")
		mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		_ = windows.SetConsoleMode(handle, mode)
	}
}
