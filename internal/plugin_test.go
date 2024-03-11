package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSiteConfigNotFound(t *testing.T) {
	p := Plugin{siteConfigs: map[string]*SiteConfig{}}

	r, err := p.getSiteConfig("test")

	var customErr = &NoSiteConfigError{}
	assert.Nil(t, r)
	assert.ErrorAs(t, err, &customErr)
}

func TestSiteConfigNoConfigs(t *testing.T) {
	p := Plugin{siteConfigs: map[string]*SiteConfig{
		"test": {},
	}}

	var customErr = &InvalidSiteConfigError{}

	r, err := p.getSiteConfig("test")
	assert.Nil(t, r)
	assert.ErrorAs(t, err, &customErr)
}
func TestSiteConfigBothConfigs(t *testing.T) {
	p := Plugin{siteConfigs: map[string]*SiteConfig{
		"test": {
			HubID: "test",
			Hubs: []HubConfig{{
				Name: "test",
			}},
		},
	}}

	var customErr = &InvalidSiteConfigError{}

	r, err := p.getSiteConfig("test")
	assert.Nil(t, r)
	assert.ErrorAs(t, err, &customErr)
}
