package mock

import (
	"bytes"

	"github.com/hvpaiva/goaoc"
)

type Manager struct {
	env           goaoc.Env
	part          string
	errSelectPart error
	errOutput     error
}

func NewBufferEnv() goaoc.Env {
	return goaoc.Env{
		Stdin:  bytes.NewBufferString(""),
		Stdout: new(bytes.Buffer),
		Args:   []string{},
	}
}

func NewManager(part string, errSelectPart, errOutput error) Manager {
	return Manager{
		env:           NewBufferEnv(),
		part:          part,
		errSelectPart: errSelectPart,
		errOutput:     errOutput,
	}
}

func (m *Manager) Read(_ string) (string, error) {
	return m.part, m.errSelectPart
}

func (m *Manager) Write(result string) error {
	if m.errOutput != nil {
		return m.errOutput
	}

	_, err := m.env.Stdout.Write([]byte(m.formatResult(result)))

	return err
}

func (m *Manager) formatResult(result string) string {
	return "The challenge result is " + result + "\n"
}

func (m *Manager) GetStdout() string {
	value, ok := m.env.Stdout.(*bytes.Buffer)
	if !ok {
		return ""
	}

	return value.String()
}
