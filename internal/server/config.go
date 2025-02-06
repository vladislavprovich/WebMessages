package server

type Config struct {
	Port      string `yaml:"SERVER_PORT" default:"8080"`
	BaseURL   string `yaml:"BASE_URL" default:"http://localhost:8080"`
	StaticDir string `yaml:"STATIC_DIR" default:"./static"`
}
