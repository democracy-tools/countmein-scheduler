package main

import (
	"fmt"

	"github.com/democracy-tools/countmein-scheduler/internal"
	"github.com/democracy-tools/countmein-scheduler/internal/slack"
)

func main() {

	sc := slack.NewClientWrapper()
	template, err := internal.Run()
	if err != nil {
		sc.Debug(fmt.Sprintf("failed to send reminder '%s' with '%s'", template, err))
		return
	}
	sc.Debug(fmt.Sprintf("reminder '%s' sent :)", template))
}
