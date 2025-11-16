package serviceapi

type Config struct {
	Address string
	BaseURL string
}

func NewConfig() (*Config, error) {
	cfg := Config{
		Address: ":8080",
	}

	//err := env.Parse(&cfg)
	//if err != nil {
	//	return nil, fmt.Errorf("error parsing .env file %w", err)
	//}

	return &cfg, nil
}
