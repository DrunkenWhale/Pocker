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
		Run()
	case "init":
		Init()
	default:
		panic("have not defined")
	}

}

func Run() {
	cmd := exec.Command(os.Args[0], append([]string{"init"}, os.Args[2])...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// this filed on be supported in linux

		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func Init() {
	syscall.Sethostname([]byte("container"))
	syscall.Mount("proc", "/proc", "proc", 0, "")
	err := syscall.Exec(os.Args[2], os.Args[2:], os.Environ())
	if err != nil {
		panic(err)
	}
	syscall.Unmount("/proc", 0)
}
