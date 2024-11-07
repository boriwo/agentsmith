package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SlackAgent struct {
}

func NewSlackAgent() *SlackAgent {
	return new(SlackAgent)
}

func (sa *SlackAgent) LaunchSlack() {
	provider, err := NewJSONSecretProvider("secrets.json")
	if err != nil {
		log.Fatalf("Error creating secret provider: %v", err)
	}
	client := slack.New(provider.GetSecret("slackOauthToken"), slack.OptionDebug(true), slack.OptionAppLevelToken(provider.GetSecret("slackAppToken")))
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		for {
			select {
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-socketClient.Events:
				switch event.Type {
				case socketmode.EventTypeEventsAPI:
					apiEvent, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPI: %v\n", event)
						continue
					}
					socketClient.Ack(*event.Request)
					err := sa.handleEventMessage(apiEvent, client)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}(ctx, client, socketClient)
	socketClient.Run()
}

func (sa *SlackAgent) handleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {
	case slackevents.CallbackEvent:
		innerEvent := event.InnerEvent
		switch evnt := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			err := sa.handleAppMentionEventToBot(evnt, client)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

func (sa *SlackAgent) handleAppMentionEventToBot(event *slackevents.AppMentionEvent, client *slack.Client) error {
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}
	text := strings.ToLower(event.Text)
	attachment := slack.Attachment{}
	if strings.Contains(text, "hello") || strings.Contains(text, "hi") {
		attachment.Text = fmt.Sprintf("Hello %s", user.RealName)
		attachment.Color = "#4af030"
	} else if strings.Contains(text, "bye") {
		attachment.Text = fmt.Sprintf("Good bye %s", user.RealName)
		attachment.Color = "#4af030"
	} else {
		attachment.Text = fmt.Sprintf("Sorry, I don't have the answers yet, %s", user.RealName)
		attachment.Color = "#4af030"
	}
	_, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

func (sa *SlackAgent) postAttachment(sp SecretProvider, pretext, text string) error {
	client := slack.New(sp.GetSecret("slackOauthToken"), slack.OptionDebug(true))
	attachment := slack.Attachment{
		Pretext: pretext,
		Text:    text,
		Color:   "4af030",
		/*Fields: []slack.AttachmentField{
			{
				Title: "Date",
				Value: time.Now().String(),
			},
		},*/
	}
	_, _, err := client.PostMessage(
		sp.GetSecret("slackChannelId"),
		slack.MsgOptionAttachments(attachment),
	)
	if err != nil {
		return err
	}
	return nil
}
