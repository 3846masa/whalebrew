package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"syscall"

	"runtime"

	"github.com/bfirsh/whalebrew/packages"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	RootCmd.AddCommand(runCommand)
}

var runCommand = &cobra.Command{
	Use:                "run PACKAGEPATH [ARGS ...]",
	Short:              "Run a package",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Help()
		}

		pkg, err := packages.LoadPackageFromPath(args[0])
		if err != nil {
			return err
		}
		dockerPath, err := exec.LookPath("docker")
		if err != nil {
			return err
		}

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		dockerArgs := []string{
			dockerPath,
			"run",
			"--interactive",
			"--rm",
			"--workdir", "/workdir",
			"-v", fmt.Sprintf("%s:/workdir", cwd),
		}
		if terminal.IsTerminal(int(os.Stdin.Fd())) {
			dockerArgs = append(dockerArgs, "--tty")
		}
		for _, volume := range pkg.Volumes {
			// special case expanding home directory
			if strings.HasPrefix(volume, "~/") {
				user, err := user.Current()
				if err != nil {
					return err
				}
				volume = user.HomeDir + volume[1:]
			}
			dockerArgs = append(dockerArgs, "-v")
			dockerArgs = append(dockerArgs, volume)
		}
		for _, envvar := range pkg.Environment {
			dockerArgs = append(dockerArgs, "-e")
			dockerArgs = append(dockerArgs, envvar)
		}
		for _, portmap := range pkg.Ports {
			dockerArgs = append(dockerArgs, "-p")
			dockerArgs = append(dockerArgs, portmap)
		}
		dockerArgs = append(dockerArgs, pkg.Image)
		dockerArgs = append(dockerArgs, args[1:]...)

		if runtime.GOOS == "windows" {
			dockerCmd := exec.Command(dockerPath, dockerArgs[1:]...)
			dockerCmd.Env = os.Environ()
			dockerCmd.Stdin = os.Stdin
			dockerCmd.Stdout = os.Stdout
			dockerCmd.Stderr = os.Stderr

			exitStatus := 1
			if err := dockerCmd.Run(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					if ws, ok := exitError.Sys().(syscall.WaitStatus); ok {
						exitStatus = ws.ExitStatus()
					}
				}
			} else {
				ws := dockerCmd.ProcessState.Sys().(syscall.WaitStatus)
				exitStatus = ws.ExitStatus()
			}
			os.Exit(exitStatus)
		}

		return syscall.Exec(dockerPath, dockerArgs, os.Environ())
	},
}
