package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/kelseyhightower/envconfig"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/slack-go/slack"
)

type config struct {
	LogLevel  string
	LogFormat string
	DryRun    bool
	Users     []string
	Groups    []string
}

type test struct {
	Package, Test, Output string
}

var cfg config

func setupLogging() {
	l, err := zerolog.ParseLevel(strings.ToLower(cfg.LogLevel))
	if err != nil {
		log.Err(err).Msg("")
		return
	}
	zerolog.SetGlobalLevel(l)

	if strings.ToLower(cfg.LogFormat) == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func readStdIn() []interface{} {
	//gojq expects []interface{} for arrays and map[]interface for single values
	//that's why this is a bit more complicated then expected, maybe this could
	//be simplified
	var (
		r []interface{}
		v interface{}
	)
	d := json.NewDecoder(os.Stdin)

	for {
		if err := d.Decode(&v); err != nil {
			break
		}
		r = append(r, v)
	}

	return r
}

func testsToReport(i []interface{}) []test {
	q, err := gojq.Parse(".[] | select(.Action==\"fail\") | select(.Test != null) | {Test,Package}")
	var r []test
	if err != nil {
		log.Err(err).Msg("can't parse tests queue")
	}
	iter := q.Run(i)
	for {
		c, ok := iter.Next()
		if !ok {
			break
		}
		t := fmt.Sprint(c.(map[string]interface{})["Test"])
		p := fmt.Sprint(c.(map[string]interface{})["Package"])
		o := testOutputsToReport(t, i)
		r = append(r, test{Package: p, Test: t, Output: o})
	}

	return r
}

func testOutputsToReport(t string, i []interface{}) string {
	q, err := gojq.Parse(fmt.Sprintf(".[] | select(.Action==\"output\") | select(.Test==\"%s\") | .Output ", t))
	if err != nil {
		log.Err(err).Msg("can't parse tests queue")
	}

	var r []string
	iter := q.Run(i)
	for {
		c, ok := iter.Next()
		if !ok {
			break
		}
		r = append(r, fmt.Sprint(c))
	}
	return strings.Join(r, "")
}

func userIdsFromEmails(c *slack.Client, emails []string) []string {
	var ids []string
	for _, e := range emails {
		user, err := c.GetUserByEmail(e)
		if err != nil {
			log.Err(err).Str("email", e).Msg("")
			continue
		}
		log.Debug().Str("ID", user.ID).Str("Fullname", user.Profile.RealName).Str("Email", user.Profile.Email).Msg("")

		ids = append(ids, user.ID)
	}

	return ids
}

func sendSlack(c *slack.Client, tests []test, usersGroups ...string) {
	messageFields := make([]slack.AttachmentField, 0)
	for _, t := range tests {
		f := slack.AttachmentField{
			Title: fmt.Sprintf("%s -> %s", t.Package, t.Test),
			Value: fmt.Sprintf("```%s```", t.Output),
		}
		messageFields = append(messageFields, f)
	}

	redHex := "#ff0000"
	attachment := slack.Attachment{
		Color:  redHex,
		Fields: messageFields,
		Text:   "Tests have failed.",
	}
	o := slack.MsgOptionAttachments(attachment)

	for _, id := range usersGroups {
		if cfg.DryRun {
			continue
		}

		if _, _, err := c.PostMessage(id, o); err != nil {
			log.Err(err).Str("id", id).Msg("can't send to message")
		}
	}

}

func main() {
	if err := envconfig.Process("notify", &cfg); err != nil {
		log.Err(err).Msg("can't read config")
	}
	setupLogging()
	log.Debug().Interface("config", cfg).Msg("")

	records := readStdIn()
	tests := testsToReport(records)
	log.Debug().Interface("Tests", tests).Msg("")

	token := "xoxb-2666384999-3765178633411-GH3pw25B54bwWuHHJsdr5oyu"
	c := slack.New(token)
	userIds := userIdsFromEmails(c, cfg.Users)

	sendSlack(c, tests, append(userIds, cfg.Groups...)...)
}
