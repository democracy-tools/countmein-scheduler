package internal

import (
	"errors"
	"fmt"
	"time"

	"github.com/democracy-tools/countmein-scheduler/internal/ds"
	"github.com/democracy-tools/countmein-scheduler/internal/whatsapp"
	"github.com/democracy-tools/go-common/env"
	"github.com/sirupsen/logrus"
)

func Run() (string, error) {

	template, err := getTemplate()
	if err != nil {
		return "", err
	}

	participantIdToPhone, err := getParticipants()
	if err != nil {
		return template, err
	}

	wac := whatsapp.NewClientWrapper()

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

func getTemplate() (string, error) {

	const location = "Asia/Jerusalem"
	isr, err := time.LoadLocation(location)
	if err != nil {
		err = fmt.Errorf("failed to load time location '%s' with '%v'", location, err)
		logrus.Error(err.Error())
		return "", err
	}

	now := time.Now().In(isr)
	nowHour := now.Hour()
	nowMin := now.In(isr).Minute()
	if (nowHour == 19 && nowMin > 57) || (nowHour == 20 && nowMin < 3) {
		return "count", nil
	}
	if nowHour == 20 && nowMin > 28 && nowMin < 33 {
		return "count2", nil
	}

	err = fmt.Errorf("invalid reminder time minute '%d:%d'", nowHour, nowMin)
	logrus.Error(err.Error())
	return "", err
}

func getParticipants() (map[string]string, error) {

	// *** datastore does not support join and group-by ***
	dsc := ds.NewClientWrapper(env.Project)

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
