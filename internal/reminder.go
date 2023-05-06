package internal

import (
	"errors"
	"fmt"

	"github.com/democracy-tools/countmein-scheduler/internal/ds"
	"github.com/democracy-tools/go-common/env"
	"github.com/democracy-tools/go-common/whatsapp"
	"github.com/sirupsen/logrus"
)

func Run(dsc ds.Client, wac whatsapp.Client) (string, error) {

	template := env.GetWhatsAppTemplate()

	participantIdToPhone, err := getParticipants(dsc)
	if err != nil {
		return template, err
	}

	var errs error
	for currUserId, currPhone := range participantIdToPhone {
		err = wac.SendReminderTemplate(template, currPhone, currUserId)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}
	if errs != nil {
		return template, errs
	}

	return template, nil
}

func getParticipants(dsc ds.Client) (map[string]string, error) {

	// *** datastore does not support join and group-by ***

	demonstration, err := ds.GetKaplanDemonstration(dsc)
	if err != nil {
		return nil, err
	}

	volunteers, err := ds.GetVolunteers(dsc, demonstration.Id)
	if err != nil {
		return nil, err
	}

	var users []ds.User
	err = dsc.GetAll(ds.KindUser, &users)
	if err != nil {
		return nil, err
	}
	userIdToPhone := make(map[string]string)
	for _, currUser := range users {
		userIdToPhone[currUser.Id] = currUser.Phone
	}

	participantIdToPhone := make(map[string]string)
	for _, currVolunteer := range volunteers {
		currPhone, ok := userIdToPhone[currVolunteer.UserId]
		if ok {
			participantIdToPhone[currVolunteer.UserId] = currPhone
		} else {
			msg := fmt.Sprintf("volunteer '%s' not found in datastore", currVolunteer.UserId)
			logrus.Error(msg)
			return nil, errors.New(msg)
		}
	}

	return participantIdToPhone, nil
}
