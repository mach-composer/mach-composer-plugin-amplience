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
		RenderTerraformProviders: state.TerraformRenderProviders,
		RenderTerraformResources: state.TerraformRenderResources,
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

	if cfg.HubID == "" && len(cfg.Hubs) == 0 {
		return nil, NewInvalidSiteConfigError("site %s must have either a hub_id or hubs set", site)
	}

	if cfg.HubID != "" && len(cfg.Hubs) > 0 {
		return nil, NewInvalidSiteConfigError("site %s cannot have both hub_id and hubs set", site)
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

func (p *Plugin) TerraformRenderProviders(site string) (string, error) {
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

func (p *Plugin) TerraformRenderResources(site string) (string, error) {
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

func (p *Plugin) RenderTerraformComponent(site string, _ string) (*schema.ComponentSchema, error) {
	cfg, err := p.getSiteConfig(site)
	if err != nil {
		return nil, err
	}

	template := `
	
		{{- if .IsMultiHub }}
			{{- range .Hubs }}
				{{ $t := printf "amplience_%s_client_id" .Name }} {{ renderProperty $t .ClientID }}
				{{ $t := printf "amplience_%s_client_secret" .Name }} {{ renderProperty $t .ClientSecret }}
				{{ $t := printf "amplience_%s_hub_id" .Name }} {{ renderProperty $t .HubID }}
			{{- end }}
		
		{{- else }}
			{{ renderProperty "amplience_client_id" .ClientID }}
			{{ renderProperty "amplience_client_secret" .ClientSecret }}
			{{ renderProperty "amplience_hub_id" .HubID }}
		{{- end }}
	`
	vars, err := helpers.RenderGoTemplate(template, cfg)
	if err != nil {
		return nil, err
	}

	var providers []string
	if cfg.IsMultiHub() {
		for _, hub := range cfg.Hubs {
			providers = append(providers, fmt.Sprintf("amplience.%s = amplience.%s", hub.Name, hub.Name))
		}
	}

	result := &schema.ComponentSchema{
		Variables: vars,
		Providers: providers,
	}
	return result, nil
}
