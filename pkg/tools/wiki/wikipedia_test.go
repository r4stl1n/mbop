package wiki

import (
	"testing"
)

func TestAdd(t *testing.T) {

	wikipediaTool := Wikipedia{}

	_, responseError := wikipediaTool.Run("dog")

	if responseError != nil {
		t.Fatal("failed")
	}
}
