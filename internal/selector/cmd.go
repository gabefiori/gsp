package selector

import (
	"bytes"
	"errors"
	"os/exec"
)

type Cmd struct {
	cmd    string
	outBuf *bytes.Buffer
	errBuf *bytes.Buffer
}

func NewCmd(cmd string) Selector {
	return &Cmd{
		cmd:    cmd,
		outBuf: new(bytes.Buffer),
		errBuf: new(bytes.Buffer),
	}
}

func (c *Cmd) Run(inputChan chan string) (string, error) {
	cmd := exec.Command(c.cmd)

	cmd.Stdout = c.outBuf
	cmd.Stderr = c.errBuf

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	go func() {
		defer stdin.Close()

		inputBuf := new(bytes.Buffer)

		for input := range inputChan {
			inputBuf.Reset()
			inputBuf.WriteString(input)
			inputBuf.WriteByte('\n')

			_, err = stdin.Write(inputBuf.Bytes())
			if err != nil {
				c.errBuf.WriteString(err.Error())
			}
		}
	}()

	_ = cmd.Wait()

	if c.errBuf.Len() > 0 {
		return "", errors.New(c.errBuf.String())
	}

	return c.outBuf.String(), nil
}
