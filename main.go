package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	fmt.Println(os.Args)
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("have not defined")
	}

}

func run() {
	cmd := exec.Command(os.Args[0], append([]string{"child"}, os.Args[2])...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// this filed on be supported in linux

		//Cloneflags: syscall.CLONE_NEWUTS|syscall.CLONE_NEWPID,

	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func child() {
	cmd := exec.Command(os.Args[2])
	syscall.Sethostname([]byte("container"))
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	syscall.Unmount("/proc", 0)
}
