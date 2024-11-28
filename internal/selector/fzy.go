package selector

import (
	"bytes"
	"errors"
	"os/exec"
)

// Fzy selector for command-line fuzzy finding.
//
// For more information, see:
// https://github.com/jhawthorn/fzy
type Fzy struct {
	outBuf *bytes.Buffer
	errBuf *bytes.Buffer
}

func NewFzy() Selector {
	return &Fzy{
		outBuf: new(bytes.Buffer),
		errBuf: new(bytes.Buffer),
	}
}

func (f *Fzy) Run(inputChan chan string) (string, error) {
	cmd := exec.Command("fzy")

	cmd.Stdout = f.outBuf
	cmd.Stderr = f.errBuf

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
				f.errBuf.WriteString(err.Error())
			}
		}
	}()

	_ = cmd.Wait()

	if f.errBuf.Len() > 0 {
		return "", errors.New(f.errBuf.String())
	}

	return f.outBuf.String(), nil
}
