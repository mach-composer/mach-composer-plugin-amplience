package internal

type GlobalConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type SiteConfig struct {
	ClientID     string      `mapstructure:"client_id"`
	ClientSecret string      `mapstructure:"client_secret"`
	HubID        string      `mapstructure:"hub_id"`
	Hubs         []HubConfig `mapstructure:"hubs"`
}

type HubConfig struct {
	Name         string `mapstructure:"name"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	HubID        string `mapstructure:"hub_id"`
}

func (s *SiteConfig) IsMultiHub() bool {
	return len(s.Hubs) > 0
}

func (s *SiteConfig) GetHubConfig(name string) *HubConfig {
	for _, hub := range s.Hubs {
		if hub.Name == name {
			return &hub
		}
	}
	return nil
}
