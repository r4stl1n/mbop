package tools

import (
	"github.com/r4stl1n/mbop/pkg/tools/wiki"
	"testing"
)

func TestAdd(t *testing.T) {

	wikipediaTool := wiki.Wikipedia{}

	_, responseError := wikipediaTool.Run("dog")

	if responseError != nil {
		t.Fatal("failed")
	}
}
