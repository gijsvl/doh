package main

import (
	"os/exec"
	"runtime"
)

func Notification(message string) { //non-intrusive notification
	if runtime.GOOS == "darwin" { // macos
		exec.Command("sh", "-c", "osascript -e '"+message+"'").Run() //notification
	}
	//TODO implement windows and linux
}
