package errors

import (
	"testing"

	"github.com/xoctopus/x/testx/bdd"
)

func TestCodeMessage_Unknown(t *testing.T) {
	bdd.From(t).Given("code is unknown", func(t bdd.T) {
		c := code(0)

		t.When("getting message", func(t bdd.T) {
			t.Then(
				"it should include SQLERROR and UNKNOWN",
				bdd.Equal(c.Message(), "[SQLERROR:0] UNKNOWN"),
			)
		})
	})
}
