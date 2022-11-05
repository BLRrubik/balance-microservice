package config

type ServerProps struct {
	Port string `yaml:"port"`
}

type DatabaseProps struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Config struct {
	ServerProps   `yaml:"server"`
	DatabaseProps `yaml:"database"`
}

func NewConfig() *Config {
	return &Config{
		ServerProps: ServerProps{
			Port: ":8080",
		},
		DatabaseProps: DatabaseProps{},
	}
}
