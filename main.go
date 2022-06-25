package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
	"syscall"
)

const usage = "Toy Project ,Just For Fun"

var (
	initCommand = cli.Command{
		Name:  "init",
		Usage: "Init container process",
		Action: func(context *cli.Context) error {
			log.Println("init come on")
			cmd := context.Args().Get(0)
			log.Printf("command %s\n", cmd)
			RunContainerInitProcess(cmd)
			return nil
		},
	}
	runCommand = cli.Command{
		Name:  "run",
		Usage: "create a container with namespace and cgroups",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "it",
				Usage: "enable tty",
			},
		},
		Action: func(context *cli.Context) error {
			if len(context.Args()) < 1 {
				return fmt.Errorf("too few argument")
			}
			cmd := context.Args().Get(0)
			tty := context.Bool("it")
			Run(tty, cmd)
			println(cmd, tty)
			return nil
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "pocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(context *cli.Context) error {
		log.SetOutput(os.Stdout)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

}

func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC,
	}
	if tty { // using std i/o
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}

func Run(tty bool, command string) {
	parent := NewParentProcess(tty, command)
	err := parent.Start()
	if err != nil {
		panic(err)
	}
	err = parent.Wait()
	if err != nil {
		panic(err)
	}
	os.Exit(-1)
}

func RunContainerInitProcess(command string) {
	log.Printf("command %s\n", command)
	defaultMountFlags := syscall.MS_NOEXEC |
		syscall.MS_NOSUID |
		syscall.MS_NODEV
	err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		panic(err)
	}
	argv := []string{command}
	err = syscall.Exec(command, argv, os.Environ())
	if err != nil {
		panic(err)
	}
}
