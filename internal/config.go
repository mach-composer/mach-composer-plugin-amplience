package internal

type AmplienceConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	HubID        string `mapstructure:"hub_id"`
}
