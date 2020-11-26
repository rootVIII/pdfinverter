package pdfinverter

import (
	"bytes"
	"os/exec"
)

// Executor runs simple shell commands.
type Executor struct {
	command bytes.Buffer
}

// SetCommand writes a command into the buffer.
func (ex *Executor) SetCommand(cmd string) {
	ex.command.WriteString(cmd)
}

// RunCommand executes the stored command before emptying the buffer.
func (ex *Executor) RunCommand() {
	cmd := exec.Command("/bin/sh", "-c", ex.command.String())
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	ex.command.Reset()
}
