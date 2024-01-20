package config

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		ConnectionString string `yaml:"connection_string"`
	} `yaml:"database"`

	JWT struct {
		SecretKey          string `yaml:"secret_key"`
		AccessTokenExpiry  int    `yaml:"access_token_expiry"`
		RefreshTokenExpiry int    `yaml:"refresh_token_expiry"`
	} `yaml:"jwt"`
}
