package bun

import (
	"io/ioutil"
	"runtime/debug"
	"testing"

	gg "github.com/danfragoso/thdwb/gg"
	"github.com/danfragoso/thdwb/ketchup"
)

func TestRenderDocument_noBody(t *testing.T) {
	html, err := ioutil.ReadFile("test-data/no-body.html")
	if err != nil {
		t.Fatalf("got unexpected error: %s", err)
	}

	doc := ketchup.ParseHTML(string(html))
	if doc == nil {
		t.Fatal("got nil document")
	}

	dctx := gg.NewContext(1024, 1024)

	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			t.Fatalf("got unexpected panic: %s. Stack: %s", err, stack)
		}
	}()

	RenderDocument(dctx, doc, false)
}
