package ctl

type Config struct {
	Version     string `mapstructure:"version"`
	Name        string `mapstructure:"name"`
	Email       string `mapstructure:"email"`
	AccessToken string `mapstructure:"access-token"`
}
