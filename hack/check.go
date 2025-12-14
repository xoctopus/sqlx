package hack

import (
	"os"
	"sync"
	"testing"
	"time"
)

var once sync.Once

func Check(t testing.TB) {
	if os.Getenv("HACK_TEST") != "true" {
		t.Skip("should depend on postgres/mysql")
	}
	once.Do(func() {
		time.Sleep(time.Second * 5) // to wait dependencies ready
	})
}
