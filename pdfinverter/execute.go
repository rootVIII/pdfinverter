package pdfinverter

import (
	"bytes"
	"os/exec"
)

// executor runs simple shell commands.
type executor struct {
	command bytes.Buffer
}

// setCommand writes a command into the buffer.
func (ex *executor) setCommand(cmd string) {
	ex.command.WriteString(cmd)
}

// runCommand executes the stored command before emptying the buffer.
func (ex *executor) runCommand() {
	cmd := exec.Command("/bin/sh", "-c", ex.command.String())
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	ex.command.Reset()
}
