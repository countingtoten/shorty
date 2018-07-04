package shorty

type Config struct {
	Address         string `env:"BIND_ADDRESS" envDefault:"localhost"`
	BindPort        int    `env:"BIND_PORT" envDefault:"3000"`
	BaseURL         string `env:"BIND_ADDRESS" envDefault:"http://localhost:3000/"`
	ShortCodeLength int    `env:"SHORT_URL_CODE_LENGTH" envDefault:"10"`
}
