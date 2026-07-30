package errors_test

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/xoctopus/x/codex"
	"github.com/xoctopus/x/testx/bdd"

	sqlerrs "github.com/xoctopus/sqlx/pkg/errors"
)

func TestIsErrNotFound(t *testing.T) {
	bdd.From(t).Given("err is NOTFOUND", func(t bdd.T) {
		err := codex.New(sqlerrs.NOTFOUND)

		t.When("checking IsErrNotFound", func(t bdd.T) {
			t.Then("it should be true", bdd.BeTrue(sqlerrs.IsErrNotFound(err)))
		})
		t.When("checking other predicates", func(t bdd.T) {
			t.Then(
				"it should not match conflict or rollback",
				bdd.BeFalse(sqlerrs.IsErrConflict(err)),
				bdd.BeFalse(sqlerrs.IsErrRollback(err)),
			)
		})
	})

	bdd.From(t).Given("err wraps NOTFOUND", func(t bdd.T) {
		err := fmt.Errorf("query failed: %w", codex.New(sqlerrs.NOTFOUND))

		t.When("checking IsErrNotFound", func(t bdd.T) {
			t.Then("it should be true", bdd.BeTrue(sqlerrs.IsErrNotFound(err)))
		})
	})

	bdd.From(t).Given("err is CONFLICT", func(t bdd.T) {
		err := codex.New(sqlerrs.CONFLICT)

		t.When("checking IsErrNotFound", func(t bdd.T) {
			t.Then("it should be false", bdd.BeFalse(sqlerrs.IsErrNotFound(err)))
		})
	})

	bdd.From(t).Given("err is a plain error", func(t bdd.T) {
		err := stderrors.New("plain")

		t.When("checking IsErrNotFound", func(t bdd.T) {
			t.Then("it should be false", bdd.BeFalse(sqlerrs.IsErrNotFound(err)))
		})
	})

	bdd.From(t).Given("err is nil", func(t bdd.T) {
		t.When("checking IsErrNotFound", func(t bdd.T) {
			t.Then("it should be false", bdd.BeFalse(sqlerrs.IsErrNotFound(nil)))
		})
	})
}

func TestIsErrConflict(t *testing.T) {
	bdd.From(t).Given("err is CONFLICT", func(t bdd.T) {
		err := codex.New(sqlerrs.CONFLICT)

		t.When("checking IsErrConflict", func(t bdd.T) {
			t.Then("it should be true", bdd.BeTrue(sqlerrs.IsErrConflict(err)))
		})
		t.When("checking other predicates", func(t bdd.T) {
			t.Then(
				"it should not match notfound or rollback",
				bdd.BeFalse(sqlerrs.IsErrNotFound(err)),
				bdd.BeFalse(sqlerrs.IsErrRollback(err)),
			)
		})
	})

	bdd.From(t).Given("err wraps CONFLICT", func(t bdd.T) {
		err := fmt.Errorf("upsert failed: %w", codex.New(sqlerrs.CONFLICT))

		t.When("checking IsErrConflict", func(t bdd.T) {
			t.Then("it should be true", bdd.BeTrue(sqlerrs.IsErrConflict(err)))
		})
	})

	bdd.From(t).Given("err is NOTFOUND", func(t bdd.T) {
		err := codex.New(sqlerrs.NOTFOUND)

		t.When("checking IsErrConflict", func(t bdd.T) {
			t.Then("it should be false", bdd.BeFalse(sqlerrs.IsErrConflict(err)))
		})
	})

	bdd.From(t).Given("err is a plain error", func(t bdd.T) {
		err := stderrors.New("plain")

		t.When("checking IsErrConflict", func(t bdd.T) {
			t.Then("it should be false", bdd.BeFalse(sqlerrs.IsErrConflict(err)))
		})
	})

	bdd.From(t).Given("err is nil", func(t bdd.T) {
		t.When("checking IsErrConflict", func(t bdd.T) {
			t.Then("it should be false", bdd.BeFalse(sqlerrs.IsErrConflict(nil)))
		})
	})
}

func TestIsErrRollback(t *testing.T) {
	bdd.From(t).Given("err is ROLLBACK", func(t bdd.T) {
		err := codex.New(sqlerrs.ROLLBACK)

		t.When("checking IsErrRollback", func(t bdd.T) {
			t.Then("it should be true", bdd.BeTrue(sqlerrs.IsErrRollback(err)))
		})
		t.When("checking other predicates", func(t bdd.T) {
			t.Then(
				"it should not match notfound or conflict",
				bdd.BeFalse(sqlerrs.IsErrNotFound(err)),
				bdd.BeFalse(sqlerrs.IsErrConflict(err)),
			)
		})
	})

	bdd.From(t).Given("err wraps ROLLBACK", func(t bdd.T) {
		err := fmt.Errorf("tx failed: %w", codex.New(sqlerrs.ROLLBACK))

		t.When("checking IsErrRollback", func(t bdd.T) {
			t.Then("it should be true", bdd.BeTrue(sqlerrs.IsErrRollback(err)))
		})
	})

	bdd.From(t).Given("err is CONFLICT", func(t bdd.T) {
		err := codex.New(sqlerrs.CONFLICT)

		t.When("checking IsErrRollback", func(t bdd.T) {
			t.Then("it should be false", bdd.BeFalse(sqlerrs.IsErrRollback(err)))
		})
	})

	bdd.From(t).Given("err is a plain error", func(t bdd.T) {
		err := stderrors.New("plain")

		t.When("checking IsErrRollback", func(t bdd.T) {
			t.Then("it should be false", bdd.BeFalse(sqlerrs.IsErrRollback(err)))
		})
	})

	bdd.From(t).Given("err is nil", func(t bdd.T) {
		t.When("checking IsErrRollback", func(t bdd.T) {
			t.Then("it should be false", bdd.BeFalse(sqlerrs.IsErrRollback(nil)))
		})
	})
}

func TestCodeMessage(t *testing.T) {
	bdd.From(t).Given("code is NOTFOUND", func(t bdd.T) {
		t.When("getting message", func(t bdd.T) {
			t.Then(
				"it should include SQLERROR and NOTFOUND",
				bdd.Equal(sqlerrs.NOTFOUND.Message(), "[SQLERROR:1] NOTFOUND"),
			)
		})
	})

	bdd.From(t).Given("code is CONFLICT", func(t bdd.T) {
		t.When("getting message", func(t bdd.T) {
			t.Then(
				"it should include SQLERROR and CONFLICT",
				bdd.Equal(sqlerrs.CONFLICT.Message(), "[SQLERROR:2] CONFLICT"),
			)
		})
	})

	bdd.From(t).Given("code is ROLLBACK", func(t bdd.T) {
		t.When("getting message", func(t bdd.T) {
			t.Then(
				"it should include SQLERROR and ROLLBACK",
				bdd.Equal(sqlerrs.ROLLBACK.Message(), "[SQLERROR:3] ROLLBACK"),
			)
		})
	})
}
