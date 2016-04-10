package web

type config struct {
	Port int `config:"port"`
}

var (
	defaultConfig = config{
		Port: 3333,
	}
)
