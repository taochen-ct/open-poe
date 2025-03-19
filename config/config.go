package config

type Configuration struct {
	App      Application `mapstructure:"internal" yaml:"internal"`
	Log      Log         `mapstructure:"log" yaml:"log"`
	Database Database    `mapstructure:"database" yaml:"database"`
	Redis    Redis       `mapstructure:"redis" yaml:"redis"`
	Limiter  Limiter     `mapstructure:"limiter" yaml:"limiter"`
}
