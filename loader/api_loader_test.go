package loader

import (
	"net/http"
	"testing"

	"api"
	stats "github.com/hellofresh/stats-go"
	"github.com/stretchr/testify/assert"
	"middleware"
	"plugin"
	"proxy"
	"router"
	"test"
	"web"
)

var tests = []struct {
	description     string
	method          string
	url             string
	headers         map[string]string
	expectedHeaders map[string]string
	expectedCode    int
}{
	{
		description: "Get example route",
		method:      "GET",
		url:         "/example",
		expectedHeaders: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
		expectedCode: http.StatusOK,
	}, {
		description: "Get invalid route",
		method:      "GET",
		url:         "/invalid-route",
		expectedHeaders: map[string]string{
			"Content-Type": "text/plain; charset=utf-8",
		},
		expectedCode: http.StatusNotFound,
	},
	{
		description: "Get one posts - strip path",
		method:      "GET",
		url:         "/posts/1",
		headers: map[string]string{
			"Host": "hellofresh.com",
		},
		expectedHeaders: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
		expectedCode: http.StatusOK,
	},
}

func TestSuccessfulLoader(t *testing.T) {
	router, err := createRegisterAndRouter()
	assert.NoError(t, err)
	ts := test.NewServer(router)
	defer ts.Close()

	for _, tc := range tests {
		res, err := ts.Do(tc.method, tc.url, tc.headers)
		assert.NoError(t, err)
		if res != nil {
			defer res.Body.Close()
		}

		for headerName, headerValue := range tc.expectedHeaders {
			assert.Equal(t, headerValue, res.Header.Get(headerName))
		}

		assert.Equal(t, tc.expectedCode, res.StatusCode, tc.description)
	}
}

func createRegisterAndRouter() (router.Router, error) {
	r := createRouter()
	r.Use(middleware.NewRecovery(web.RecoveryHandler).Handler)

	register := proxy.NewRegister(r, createProxy())
	proxyRepo, err := createProxyRepo()
	if err != nil {
		return nil, err
	}

	pluginLoader := plugin.NewLoader()
	loader := NewAPILoader(register, pluginLoader)
	loader.LoadDefinitions(proxyRepo)

	return r, nil
}

func createProxyRepo() (api.Repository, error) {
	return api.NewFileSystemRepository("../../examples/apis")
}

func createRouter() router.Router {
	return router.NewChiRouter()
}

func createProxy() *proxy.Proxy {
	return proxy.WithParams(proxy.Params{
		StatsClient: stats.NewStatsdClient("", ""),
	})
}
