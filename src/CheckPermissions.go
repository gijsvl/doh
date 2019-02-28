package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func CheckPermissions() { //CHECK PERMISSIONS
	if os.Geteuid() != 0 { //check if not root
		if runtime.GOOS == "darwin" { // macos
			command := "osascript -e 'do shell script \"" + os.Args[0] + "\" with prompt \"DoH needs system rights\" with administrator privileges'"
			err := exec.Command("sh", "-c", command).Run()
			CheckError(err)
			os.Exit(0)
		} else if runtime.GOOS == "windows" { //windows
			//TODO test if this works?
			err := exec.Command("/runas", "/profile", "/user:administrator", os.Args[0]).Run()
			CheckError(err)
			os.Exit(0)
		} else { //linux
			fmt.Println("Please run this program as root")
		}
		os.Exit(126)
	}
}
