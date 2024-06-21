package internal

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func normalize(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), ""))
}

func TestGetSiteConfigNotFound(t *testing.T) {
	p := Plugin{siteConfigs: map[string]*SiteConfig{}}

	r, err := p.getSiteConfig("test")

	var customErr = &NoSiteConfigError{}
	assert.Nil(t, r)
	assert.ErrorAs(t, err, &customErr)
}

func TestSetSiteConfigNoConfigs(t *testing.T) {
	p := Plugin{}

	var customErr = &InvalidSiteConfigError{}

	err := p.SetSiteConfig("test", map[string]any{})
	assert.ErrorAs(t, err, &customErr)
}

func TestSetSiteConfigBothConfigs(t *testing.T) {
	p := Plugin{}

	var customErr = &InvalidSiteConfigError{}

	err := p.SetSiteConfig("test", map[string]any{
		"client_id":     "test",
		"client_secret": "test",
		"hub_id":        "test",
		"hubs": []map[string]any{
			{
				"name":      "test",
				"alias":     "test",
				"client_id": "test",
				"hub_id":    "test",
			},
		},
	})
	assert.ErrorAs(t, err, &customErr)
}

func TestRenderTerraformResourcesSiteNotFound(t *testing.T) {
	p := Plugin{siteConfigs: map[string]*SiteConfig{
		"site": {},
	}}

	_, err := p.RenderTerraformResources("site2")
	assert.Error(t, err)
}

func TestRenderTerraformResourcesSingleHub(t *testing.T) {
	p := Plugin{siteConfigs: map[string]*SiteConfig{
		"test": {ClientID: "test", ClientSecret: "test", HubID: "test"},
	}}

	r, err := p.RenderTerraformResources("test")

	assert.Nil(t, err)
	assert.Equal(t, normalize(`
		provider "amplience" {
			client_id     = "test"
			client_secret = "test"
			hub_id        = "test"
		}
	`), normalize(r))
}

func TestRenderTerraformResourcesMultiHub(t *testing.T) {
	p := Plugin{siteConfigs: map[string]*SiteConfig{
		"test": {Hubs: []HubConfig{
			{Name: "test1", ClientID: "test1", ClientSecret: "test1", HubID: "test1"},
			{Name: "test2", ClientID: "test2", ClientSecret: "test2", HubID: "test2"},
		}},
	}}

	r, err := p.RenderTerraformResources("test")

	assert.Nil(t, err)
	assert.Equal(t, normalize(`
		provider "amplience" {
			alias         = "test1"
			client_id     = "test1"
			client_secret = "test1"
			hub_id        = "test1"
		}
		
		provider "amplience" {
			alias         = "test2"
			client_id     = "test2"
			client_secret = "test2"
			hub_id        = "test2"
		}
	
	`), normalize(r))
}

func TestRenderTerraformComponentSiteNotFound(t *testing.T) {
	p := Plugin{siteConfigs: map[string]*SiteConfig{
		"site": {},
	}}

	_, err := p.RenderTerraformComponent("site2", "test2")
	assert.Error(t, err)
}

func TestRenderTerraformComponentSimple(t *testing.T) {
	p := Plugin{siteConfigs: map[string]*SiteConfig{
		"site": {ClientID: "test", ClientSecret: "test", HubID: "test"},
	}}

	r, err := p.RenderTerraformComponent("site", "component")

	assert.Nil(t, err)
	assert.Equal(t, normalize(`
		amplience_client_id = "test"
		amplience_client_secret = "test"
		amplience_hub_id = "test"
	`), normalize(r.Variables))
	assert.Equal(t, 0, len(r.Providers))
}

func TestRenderTerraformComponentMultiHub(t *testing.T) {
	p := Plugin{siteConfigs: map[string]*SiteConfig{
		"site": {Hubs: []HubConfig{
			{Name: "test1", ClientID: "test1", ClientSecret: "test1", HubID: "test1"},
			{Name: "test2", ClientID: "test2", ClientSecret: "test2", HubID: "test2"},
		}},
	}}

	r1, err := p.RenderTerraformComponent("site", "test1")
	assert.Nil(t, err)
	assert.Equal(t, normalize(`
		amplience_client_id = "test1"
		amplience_client_secret = "test1"
		amplience_hub_id = "test1"
	`), normalize(r1.Variables))
	assert.Equal(t, 1, len(r1.Providers))
	assert.Equal(t, "amplience = amplience.test1", r1.Providers[0])

	r2, err := p.RenderTerraformComponent("site", "test2")
	assert.Nil(t, err)
	assert.Equal(t, normalize(`
		amplience_client_id = "test2"
		amplience_client_secret = "test2"
		amplience_hub_id = "test2"
	`), normalize(r2.Variables))
	assert.Equal(t, 1, len(r2.Providers))
	assert.Equal(t, "amplience = amplience.test2", r2.Providers[0])
}

func TestRenderTerraformComponentHubNotFound(t *testing.T) {
	p := Plugin{siteConfigs: map[string]*SiteConfig{
		"site": {Hubs: []HubConfig{
			{Name: "test1", ClientID: "test1", ClientSecret: "test1", HubID: "test1"},
		}},
	}}

	var customErr = &NoHubConfigError{}

	_, err := p.RenderTerraformComponent("site", "test2")
	assert.ErrorAs(t, err, &customErr)
}

func TestRenderTerraformProviders(t *testing.T) {
	p := Plugin{
		provider: "1.0.0",
		siteConfigs: map[string]*SiteConfig{
			"site": {ClientID: "test1", ClientSecret: "test1", HubID: "test1"},
		},
	}

	r, err := p.RenderTerraformProviders("site")
	assert.Nil(t, err)
	assert.Equal(t, normalize(`
		amplience = {
			source = "labd/amplience"
			version = "~>1.0.0"
		}
	`), normalize(r))
}

func TestRenderTerraformProvidersSiteNotFound(t *testing.T) {
	p := Plugin{siteConfigs: map[string]*SiteConfig{
		"site": {},
	}}

	_, err := p.RenderTerraformProviders("site2")
	assert.Error(t, err)
}
