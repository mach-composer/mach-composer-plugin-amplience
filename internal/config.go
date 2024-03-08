package internal

type AmplienceConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	HubID        string `mapstructure:"hub_id"`
	ExtraHubs    []struct {
		Name         string `mapstructure:"name"`
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		HubID        string `mapstructure:"hub_id"`
	} `mapstructure:"extra_hubs"`
}
