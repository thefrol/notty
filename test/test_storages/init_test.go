package storages_test

import (
	"fmt"
	"os"
	"testing"

	_ "gitlab.com/thefrol/notty/internal/migrations"
)

const (
	EnvTestConnection = "NOTTY_TEST_DB"
	EnvTestString     = "NOTTY_TEST_NOSKIP"
)

// Локальные тесты можно пропускать на этом шаге, но
// на гитлабе надо запретить такую штуку при установке
// NOTTY_TEST_NOSKIP=1
var NoSkip bool

var TestDSN string

func init() {
	if os.Getenv(EnvTestString) == "1" {
		NoSkip = true
	}

	TestDSN = os.Getenv(EnvTestConnection)

	fmt.Println(NoSkip, TestDSN)
}

func Test_nothin(t *testing.T) {

}
