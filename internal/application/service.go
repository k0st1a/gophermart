package application

func Run() error {
	cfg, err := collectConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	return nil
}
