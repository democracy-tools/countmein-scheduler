package main

import (
	"fmt"

	"github.com/democracy-tools/countmein-scheduler/internal"
	"github.com/democracy-tools/countmein-scheduler/internal/ds"
	"github.com/democracy-tools/go-common/env"
	"github.com/democracy-tools/go-common/slack"
	"github.com/democracy-tools/go-common/whatsapp"
)

func main() {

	sc := slack.NewClientWrapper()
	template, err := internal.Run(ds.NewClientWrapper(env.Project), whatsapp.NewClientWrapper())
	if err != nil {
		sc.Debug(fmt.Sprintf("failed to send reminder '%s' with '%s'", template, err))
		return
	}
	sc.Debug(fmt.Sprintf("reminder '%s' sent :)", template))
}
