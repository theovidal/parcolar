package pronote

import (
	"fmt"
	"github.com/theovidal/parcoursupbot/lib"
	"testing"
)

func TestGetHomework(t *testing.T) {
	lib.LoadEnv("../.env")
	lib.OpenCache()

	response, err := GetHomework()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(response.Data.Homeworks[0].String())
}
