package config

import (
	"cef/pkg/external/aegis"
	"testing"
)

func TestLoader_GetBrowserConfigLoader(t *testing.T) {
	var l Loader
	if err := l.LoadExternalConfig(); err != nil {
		t.Fatal(err)
	}
	aegis.SetDefault(aegis.NewAegisClient(l.ExternalConfig.AegisAddr.Mode))
	browserConfigLoader := l.GetBrowserConfigLoader()
	conf := browserConfigLoader("test")
	t.Log(conf)
}
