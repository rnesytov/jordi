package config

type Config struct {
	Target   string
	Method   string
	Insecure bool
}

func New(target, method string, insecure bool) Config {
	return Config{Target: target, Method: method, Insecure: insecure}
}

func (c Config) Validate() error {
	return nil
}
