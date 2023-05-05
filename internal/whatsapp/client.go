package whatsapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/democracy-tools/countmein-scheduler/internal/env"
	"github.com/sirupsen/logrus"
)

type Client interface {
	SendReminderTemplate(template string, phone string, userId string) error
}

type ClientWrapper struct {
	auth string
	from string
}

func NewClientWrapper() Client {
	return &ClientWrapper{
		auth: fmt.Sprintf("Bearer %s", env.GetWhatAppToken()),
		from: env.GetWhatsAppFromPhone(),
	}
}

func (c *ClientWrapper) SendReminderTemplate(template string, to string, userId string) error {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(newTemplate(template, to, fmt.Sprintf("?user-id=%s", userId), nil))
	if err != nil {
		err = fmt.Errorf("failed to encode whatsapp reminder message with '%v' phone '%s'", err, to)
		logrus.Error(err.Error())
		return err
	}

	return send(c.from, to, &buf, c.auth)
}

func send(from string, to string, body io.Reader, auth string) error {

	r, err := http.NewRequest(http.MethodPost, getMessageUrl(from), body)
	if err != nil {
		logrus.Errorf("failed to create HTTP request for sending a whatsapp message to '%s' with '%v'", to, err)
		return err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", auth)

	client := http.Client{}
	response, err := client.Do(r)
	if err != nil {
		logrus.Errorf("failed to send whatsapp message to '%s' with '%v'", to, err)
		return err
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		msg := fmt.Sprintf("failed to send whatsapp message to '%s' with '%s'", to, response.Status)
		logrus.Info(msg)
		return errors.New(msg)
	}

	return nil
}

func newTemplate(name string, to string, buttonUrlParam string, bodyTextParams []string) TemplateMessageRequest {

	res := TemplateMessageRequest{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: Template{
			Name: name,
			Language: Language{
				Policy: "deterministic",
				Code:   "he",
			},
		},
	}

	if buttonUrlParam != "" {
		res.Template.Components = []Component{{
			Type:    "button",
			SubType: "url",
			Index:   "0",
			Parameters: []Parameter{{
				Type: "text",
				Text: buttonUrlParam,
			}},
		}}
	}

	if len(bodyTextParams) > 0 {
		var params []Parameter
		for _, currParam := range bodyTextParams {
			params = append(params, Parameter{
				Type: "text",
				Text: currParam,
			})
		}
		res.Template.Components = append(res.Template.Components, Component{
			Type:       "body",
			Parameters: params})
	}

	return res
}

func getMessageUrl(from string) string {

	return fmt.Sprintf("https://graph.facebook.com/v16.0/%s/messages", from)
}