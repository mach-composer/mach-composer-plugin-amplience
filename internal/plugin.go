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
		siteConfigs: map[string]*AmplienceConfig{},
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
	globalConfig *AmplienceConfig
	siteConfigs  map[string]*AmplienceConfig
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
	cfg := AmplienceConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.globalConfig = &cfg
	p.enabled = true
	return nil
}

func (p *Plugin) SetSiteConfig(site string, data map[string]any) error {
	cfg := AmplienceConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.siteConfigs[site] = &cfg
	p.enabled = true
	return nil
}

func (p *Plugin) getSiteConfig(site string) *AmplienceConfig {
	result := &AmplienceConfig{}
	if p.globalConfig != nil {
		result.ClientID = p.globalConfig.ClientID
		result.ClientSecret = p.globalConfig.ClientSecret
		result.HubID = p.globalConfig.HubID
	}

	cfg, ok := p.siteConfigs[site]
	if ok {
		if cfg.ClientID != "" {
			result.ClientID = cfg.ClientID
		}
		if cfg.ClientSecret != "" {
			result.ClientSecret = cfg.ClientSecret
		}
		if cfg.HubID != "" {
			result.HubID = cfg.HubID
		}
		result.ExtraHubs = cfg.ExtraHubs
	}

	if result.ClientID == "" {
		return nil
	}
	return result
}

func (p *Plugin) TerraformRenderProviders(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	result := fmt.Sprintf(`
	amplience = {
		source = "labd/amplience"
		version = "%s"
	}`, helpers.VersionConstraint(p.provider))
	return result, nil
}

func (p *Plugin) TerraformRenderResources(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	template := `
		provider "amplience" {
			{{ renderProperty "client_id" .ClientID }}
			{{ renderProperty "client_secret" .ClientSecret }}
			{{ renderProperty "hub_id" .HubID }}
		}
	
		{{range .ExtraHubs}}
	    provider "amplience" {
			{{ renderProperty "alias" .Name }}
			{{ renderProperty "client_id" .ClientID }}
			{{ renderProperty "client_secret" .ClientSecret }}
			{{ renderProperty "hub_id" .HubID }}
		}
		{{end}}
	`
	return helpers.RenderGoTemplate(template, cfg)
}

func (p *Plugin) RenderTerraformComponent(site string, _ string) (*schema.ComponentSchema, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return nil, nil
	}

	template := `
		{{ renderProperty "amplience_client_id" .ClientID }}
		{{ renderProperty "amplience_client_secret" .ClientSecret }}
		{{ renderProperty "amplience_hub_id" .HubID }}
	
		{{- range .ExtraHubs }}
		{{- with $key := printf "amplience_%s_client_id" .Name }}
			{{ renderProperty . $.ClientID }}
		{{- end }}
		{{- with $key := printf "amplience_%s_client_secret" .Name }}
			{{ renderProperty . $.ClientSecret }}
		{{- end }}
		{{- with $key := printf "amplience_%s_hub_id" .Name }}
			{{ renderProperty . $.HubID }}
		{{- end }}
		{{- end }}
	`
	vars, err := helpers.RenderGoTemplate(template, cfg)
	if err != nil {
		return nil, err
	}

	var providers = []string{"amplience = amplience"}

	for _, hub := range cfg.ExtraHubs {
		providers = append(providers, fmt.Sprintf("amplience.%s = amplience.%s", hub.Name, hub.Name))
	}

	result := &schema.ComponentSchema{
		Variables: vars,
		Providers: providers,
	}
	return result, nil
}
