package config

type Config struct {
	Target string
}

func New(target string) Config {
	return Config{Target: target}
}

func (c Config) Validate() error {
	return nil
}
