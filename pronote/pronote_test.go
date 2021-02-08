package pronote

import (
	"fmt"
	"testing"

	"github.com/theovidal/bacbot/lib"
	"github.com/theovidal/bacbot/pronote/api"
)

func TestGetHomework(t *testing.T) {
	lib.LoadEnv("../.env")
	lib.OpenCache()

	response, err := api.GetHomework()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(response.Homeworks[0].String())
}

func TestGetTimetable(t *testing.T) {
	lib.LoadEnv("../.env")
	lib.OpenCache()

	response, err := api.GetTimetable(false)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(response.Timetable[0].String())
}
