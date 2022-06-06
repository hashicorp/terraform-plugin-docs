package provider

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_resourceName(t *testing.T) {
	cases := []struct {
		name                 string
		providerShortName    string
		templateFileName     string
		expectedResourceName string
	}{
		{
			"provider short name same as template file name",
			"http",
			"http.md.tmpl",
			"http",
		},
		{
			"provider short name different to template file name",
			"tls",
			"cert_request.md.tmpl",
			"tls_cert_request",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actualResourceName := resourceName(c.providerShortName, c.templateFileName)
			if !cmp.Equal(c.expectedResourceName, actualResourceName) {
				t.Errorf("expected: %s, got: %s", c.expectedResourceName, actualResourceName)
			}
		})
	}
}
