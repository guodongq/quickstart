package provider

type Provider interface {
	Init() error
	Close() error
}

type RunProvider interface {
	Provider

	Run() error
	IsRunning() bool
}

type AbstractProvider struct {
	Provider
}

func (p *AbstractProvider) Init() error {
	return nil
}

func (p *AbstractProvider) Close() error {
	return nil
}

type AbstractRunProvider struct {
	RunProvider

	running bool
}

func (p *AbstractRunProvider) Init() error {
	return nil
}

func (p *AbstractRunProvider) SetRunning(running bool) {
	p.running = running
}

func (p *AbstractRunProvider) IsRunning() bool {
	return p.running
}

func (p *AbstractRunProvider) Close() error {
	p.SetRunning(false)
	return nil
}
