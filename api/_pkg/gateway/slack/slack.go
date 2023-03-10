// Copyright 2023 ainsley.dev LTD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slack

import (
	"context"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type (
	// Sender defines the method to send messages via the Slack API.
	Sender interface {
		// Send takes a message subject and a message body and sends to the set channel.
		// A Client app with the chat:write.public and chat:write permissions must
		// be installed to the workspace.
		// See: https://api.slack.com/
		Send(ctx context.Context, channelID, subject, message string) error
	}
	// Client implements the notifier interface to send Client.
	// messages with a given message.
	Client struct {
		slackSendFunc slackSendFn
	}
	// sendSlackFunc is the function used for sending to
	// a Client channel, it's stubbed for testing.
	slackSendFn func(ctx context.Context, channelID string, options ...slack.MsgOption) (string, string, error)
)

var (
	// Channels represents the key value pairs of thread
	// channels for Slack.
	Channels = struct {
		Contact string
	}{
		Contact: "#website-contact",
	}
)

// New creates a new Client notifier.
// For more information about slack API token:
//
//	-> https://pkg.go.dev/github.com/slack-go/slack#New
func New(token string) *Client {
	return &Client{
		slackSendFunc: slack.New(token).PostMessageContext,
	}
}

func (c *Client) Send(ctx context.Context, channelID, subject, message string) error {
	// Create the Client attachment that we will send to the channel.
	attachment := slack.Attachment{
		Pretext: subject,
		Text:    message,
	}
	id, timestamp, err := c.slackSendFunc(ctx, channelID, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return errors.Wrapf(err,
			"failed to send message to Slack channel '%s' at time '%s'",
			id,
			timestamp,
		)
	}
	return nil
}
