package sauce

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	structs "thdwb/structs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`OK!`))
	}))

	defer server.Close()

	url, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	resource := GetResource(url, nil)
	var testResorce *structs.Resource

	assert.IsType(t, testResorce, resource, "Expecting: Resource Struct")
	assert.Equal(t, "OK!", resource.Body, "Expecting: OK!")
}

func TestGetImage(t *testing.T) {
	testCases := []struct {
		name           string
		url            string
		expectResponse string
	}{
		{"base64 encoded", "", ""},
		{"regular URL", "", ""},
		{"broken URL", "", ""},
		{"empty URL", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := url.Parse(tc.url)
			require.NoError(t, err)

			var data []byte

			assert.NotPanics(t, func() {
				data = GetImage(u)
			})

			assert.NotEqual(t, 0, len(data))
		})
	}
}
