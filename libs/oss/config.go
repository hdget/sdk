package oss

type Config struct {
	Region       string `mapstructure:"region"`
	Domain       string `mapstructure:"domain"`
	Bucket       string `mapstructure:"bucket"`
	AccessKey    string `mapstructure:"access_key"`
	AccessSecret string `mapstructure:"access_secret"`
}
