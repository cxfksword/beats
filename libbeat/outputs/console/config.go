package console

type config struct {
	Pretty bool   `config:"pretty"`
	Query  string `config:"query"`
}

var (
	defaultConfig = config{
		Pretty: false,
	}
)
