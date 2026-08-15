package catalog

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRemoteProviderResolvesCasaOSCompose(t *testing.T) {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/index.json":
			_, _ = fmt.Fprintf(w, `{"base_url":%q,"apps":[{"id":"org.example.demo","title":"Demo","tagline":"Demo app","category":"Utilities","compose_url":"/demo/docker-compose.yml"}]}`, server.URL)
		case "/demo/docker-compose.yml":
			_, _ = w.Write([]byte("services:\n  demo:\n    image: example/demo:latest\n"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := RemoteProvider{URL: server.URL + "/index.json", Client: server.Client()}
	apps, err := provider.ListApps(context.Background())
	if err != nil || len(apps) != 1 || !apps[0].Installable {
		t.Fatalf("list = %+v, err = %v", apps, err)
	}
	app, err := provider.GetApp("org.example.demo")
	if err != nil || app.Compose == "" {
		t.Fatalf("app = %+v, err = %v", app, err)
	}
}
