package site

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func discoveryFixture(t *testing.T, setting string) (*Site, http.Handler) {
	t.Helper()
	files := cloneMapFS(fixture)
	config := files["site.json"]
	config.Data = []byte(strings.Replace(string(config.Data), `"nav":`, setting+`"nav":`, 1))
	s, h, err := NewServerHandler(files, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	if err := Validate(s); err != nil {
		t.Fatal(err)
	}
	return s, h
}

func discoveryResponse(h http.Handler, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(method, path, nil))
	return w
}

func TestMCPDiscoveryOptOut(t *testing.T) {
	for _, tc := range []struct {
		name, setting string
		enabled       bool
	}{
		{"default", "", true}, {"off", `"mcp":"off",`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s, h := discoveryFixture(t, tc.setting)
			if strings.Contains(s.RenderAgentsTxt(), "\nmcp: ") != tc.enabled || !strings.Contains(s.RenderAgentsTxt(), "bashy site attach") {
				t.Fatal("incorrect agents.txt discovery")
			}
			for _, path := range []string{"/", "/index.md", "/index.json", "/agents.txt", "/AGENTS.md", "/skill.md", "/llms.txt", "/humans.txt", "/robots.txt", "/.well-known/api-catalog", "/feed.xml", "/missing"} {
				w := discoveryResponse(h, "GET", path)
				if strings.Contains(strings.ToLower(w.Body.String()), "mcp") != tc.enabled || strings.Contains(w.Header().Get("Link"), "service-desc") != tc.enabled || strings.Contains(w.Header().Get("Bashy-Shaped"), "mcp=") != tc.enabled {
					t.Errorf("%s: MCP discovery disagrees with config", path)
				}
			}
			for _, path := range []string{"/.well-known/mcp/server-card.json", "/.well-known/oauth-protected-resource"} {
				for _, method := range []string{"GET", "HEAD"} {
					want := http.StatusOK
					if !tc.enabled {
						want = http.StatusNotFound
					}
					if w := discoveryResponse(h, method, path); w.Code != want {
						t.Errorf("%s %s: got %d want %d", method, path, w.Code, want)
					}
				}
			}
		})
	}
}

func TestSiteAttachOptOut(t *testing.T) {
	for _, tc := range []struct {
		name, setting string
		enabled       bool
	}{
		{"default", "", true}, {"enabled", `"site_attach":true,`, true}, {"disabled", `"site_attach":false,`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s, h := discoveryFixture(t, tc.setting)
			if strings.Contains(s.AgentRecipe(s.homePage()), "bashy site attach") != tc.enabled || !strings.Contains(s.RenderAgentsTxt(), "\nmcp: ") {
				t.Fatal("incorrect agent recipe or MCP discovery")
			}
			for _, path := range []string{"/", "/index.md", "/index.json", "/agents.txt", "/AGENTS.md", "/llms.txt", "/feed.xml"} {
				w := discoveryResponse(h, "GET", path)
				if w.Code != http.StatusOK || strings.Contains(w.Body.String(), "bashy site attach") != tc.enabled {
					t.Errorf("%s: site attach discovery disagrees with config", path)
				}
			}
		})
	}
}
