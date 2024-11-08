package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SlackAgent struct {
	client         *slack.Client
	secretProvider SecretProvider
	answerProvider AnswerProvider
	sessionMgr     SessionManager
}

func NewSlackAgent(secretProvider SecretProvider, answerProvider AnswerProvider, sessionManager SessionManager) *SlackAgent {
	sa := SlackAgent{
		secretProvider: secretProvider,
		answerProvider: answerProvider,
		sessionMgr:     sessionManager,
	}
	return &sa
}

func (sa *SlackAgent) LaunchAgent() {
	sa.client = slack.New(sa.secretProvider.GetSecret("slackOauthToken"), slack.OptionDebug(true), slack.OptionAppLevelToken(sa.secretProvider.GetSecret("slackAppToken")))
	socketClient := socketmode.New(
		sa.client,
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
	}(ctx, sa.client, socketClient)
	sa.postAttachment("Bot Message", "agent launched")
	socketClient.Run()
	sa.postAttachment("Bot Message", "agent stopped")
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
	slackUser, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}
	user := NewUser(slackUser.ID, slackUser.Name, slackUser.RealName)
	session := sa.sessionMgr.GetSession(user)
	question := NewQuestion(event.Text)
	answers := sa.answerProvider.GetAnswers(session, question)
	for _, a := range answers {
		attachment := slack.Attachment{}
		attachment.Text = a.Text
		attachment.Color = "#4af030"
		_, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
		if err != nil {
			return fmt.Errorf("failed to post message: %w", err)
		}
	}
	return nil
}

func (sa *SlackAgent) postAttachment(pretext, text string) error {
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
	_, _, err := sa.client.PostMessage(
		sa.secretProvider.GetSecret("slackChannelId"),
		slack.MsgOptionAttachments(attachment),
	)
	if err != nil {
		return err
	}
	return nil
}
