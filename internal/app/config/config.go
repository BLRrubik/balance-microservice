package config

type ServerProps struct {
	Port string `yaml:"port"`
}

type LoggerProps struct {
	Level string `yaml:"level"`
}

type DatabaseProps struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type CookieProps struct {
	Secret string `yaml:"secret"`
}

type Config struct {
	ServerProps   `yaml:"server"`
	LoggerProps   `yaml:"logger"`
	DatabaseProps `yaml:"database"`
}

func NewConfig() *Config {
	return &Config{
		ServerProps: ServerProps{
			Port: ":8080",
		},
		LoggerProps: LoggerProps{
			Level: "debug",
		},
		DatabaseProps: DatabaseProps{},
	}
}
