package sauce

import (
	"encoding/base64"
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
	// This is a base64-encoded 1x1 PNG containing a single opaque black pixel
	blackPixel := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoded, err := base64.RawStdEncoding.DecodeString(blackPixel)
		require.NoError(t, err)

		w.Write(decoded)
	}))
	defer srv.Close()

	testCases := []struct {
		name      string
		url       string
		expectLen int
	}{
		{"base64 encoded", "data:image/png;base64," + blackPixel, 66},
		{"malformed base64 encoded 1", "data:image/png;base64:" + blackPixel, 0},
		{"malformed base64 encoded 2", "data:image/png;base64", 0},
		{"short URL", "http://foo", 0},
		{"regular URL", srv.URL + "/this-is-a-relatively-long-url", 0},
		{"not an HTTP URL", "//foo-bar-fnord-asdf-bla-bla", 0},
		{"empty URL", "", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := url.Parse(tc.url)
			require.NoError(t, err)

			var data []byte

			assert.NotPanics(t, func() {
				data = GetImage(u)
			})

			assert.Len(t, data, tc.expectLen)
		})
	}
}
