package internal

import (
	"fmt"
	"github.com/mach-composer/mach-composer-plugin-helpers/helpers"
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/mitchellh/mapstructure"
)

func NewAmpliencePlugin() schema.MachComposerPlugin {
	state := &Plugin{
		provider:    "0.3.7",
		siteConfigs: map[string]*SiteConfig{},
	}

	return plugin.NewPlugin(&schema.PluginSchema{
		Identifier: "amplience",

		Configure: state.Configure,
		IsEnabled: func() bool { return state.enabled },

		GetValidationSchema: state.GetValidationSchema,

		// Config
		SetGlobalConfig: state.SetGlobalConfig,
		SetSiteConfig:   state.SetSiteConfig,

		// Renders
		RenderTerraformProviders: state.RenderTerraformProviders,
		RenderTerraformResources: state.RenderTerraformResources,
		RenderTerraformComponent: state.RenderTerraformComponent,
	})
}

type Plugin struct {
	environment  string
	provider     string
	globalConfig *GlobalConfig
	siteConfigs  map[string]*SiteConfig
	enabled      bool
}

func (p *Plugin) Configure(environment string, provider string) error {
	p.environment = environment
	if provider != "" {
		p.provider = provider
	}
	return nil
}

func (p *Plugin) GetValidationSchema() (*schema.ValidationSchema, error) {
	result := getSchema()
	return result, nil
}

func (p *Plugin) SetGlobalConfig(data map[string]any) error {
	cfg := GlobalConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.globalConfig = &cfg
	p.enabled = true
	return nil
}

func (p *Plugin) SetSiteConfig(site string, data map[string]any) error {
	cfg := SiteConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}

	if cfg.HubID == "" && len(cfg.Hubs) == 0 {
		return NewInvalidSiteConfigError("site %s must have either a hub_id or hubs set", site)
	}

	if cfg.HubID != "" && len(cfg.Hubs) > 0 {
		return NewInvalidSiteConfigError("site %s cannot have both hub_id and hubs set", site)
	}

	p.siteConfigs[site] = &cfg
	p.enabled = true
	return nil
}

func (p *Plugin) getSiteConfig(site string) (*SiteConfig, error) {
	result := &SiteConfig{}
	if p.globalConfig != nil {
		result.ClientID = p.globalConfig.ClientID
		result.ClientSecret = p.globalConfig.ClientSecret
	}

	cfg, ok := p.siteConfigs[site]
	if !ok {
		return nil, NewNoSiteConfigError("site %s not found", site)
	}

	if cfg.ClientID != "" {
		result.ClientID = cfg.ClientID
	}
	if cfg.ClientSecret != "" {
		result.ClientSecret = cfg.ClientSecret
	}
	if cfg.HubID != "" {
		result.HubID = cfg.HubID
	}

	result.Hubs = cfg.Hubs

	return result, nil
}

func (p *Plugin) RenderTerraformProviders(site string) (string, error) {
	_, err := p.getSiteConfig(site)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf(`
	amplience = {
		source = "labd/amplience"
		version = "%s"
	}`, helpers.VersionConstraint(p.provider))
	return result, nil
}

func (p *Plugin) RenderTerraformResources(site string) (string, error) {
	cfg, err := p.getSiteConfig(site)
	if err != nil {
		return "", err
	}

	template := `
		{{- if .IsMultiHub }}
		{{- range .Hubs }}
	    provider "amplience" {
			{{ renderProperty "alias" .Name }}
			{{ renderProperty "client_id" .ClientID }}
			{{ renderProperty "client_secret" .ClientSecret }}
			{{ renderProperty "hub_id" .HubID }}
		}
		{{- end }}
		{{- else }}
		provider "amplience" {
			{{ renderProperty "client_id" .ClientID }}
			{{ renderProperty "client_secret" .ClientSecret }}
			{{ renderProperty "hub_id" .HubID }}
		}
		{{- end }}
	`
	return helpers.RenderGoTemplate(template, cfg)
}

type tplConfig struct {
	ClientID     string
	ClientSecret string
	HubID        string
}

func (p *Plugin) RenderTerraformComponent(site string, component string) (*schema.ComponentSchema, error) {
	siteConfig, err := p.getSiteConfig(site)
	if err != nil {
		return nil, err
	}

	tplData := tplConfig{
		ClientID:     siteConfig.ClientID,
		ClientSecret: siteConfig.ClientSecret,
		HubID:        siteConfig.HubID,
	}

	var providers []string
	if siteConfig.IsMultiHub() {
		hubCfg := siteConfig.GetHubConfig(component)
		if hubCfg == nil {
			return nil, NewNoHubConfigError("hub %s not found in site %s. A hub must exist with the same name as the component", component, site)
		}

		tplData.ClientID = hubCfg.ClientID
		tplData.ClientSecret = hubCfg.ClientSecret
		tplData.HubID = hubCfg.HubID

		providers = append(providers, fmt.Sprintf("amplience = amplience.%s", component))
	}

	template := `
		{{ renderProperty "amplience_client_id" .ClientID }}
		{{ renderProperty "amplience_client_secret" .ClientSecret }}
		{{ renderProperty "amplience_hub_id" .HubID }}
	`
	vars, err := helpers.RenderGoTemplate(template, tplData)
	if err != nil {
		return nil, err
	}

	result := &schema.ComponentSchema{
		Variables: vars,
		Providers: providers,
	}
	return result, nil
}
