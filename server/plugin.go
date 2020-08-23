package main

import (
	"fmt"
	"strings"
	"sync"

	"math/rand"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin
	configurationLock sync.RWMutex
	configuration     *configuration
}

const spongeBob string = ":spongebob:"
const spongeBobLen int = len(spongeBob)

func (p *Plugin) OnActivate() (err error) {
	// args.Command contains the full command string entered
	err = p.API.RegisterCommand(&model.Command{
		Trigger:          "spongebob",
		DisplayName:      "Mocking Spongebob Text",
		Description:      "iT aPpEaRs LikE ThiS :spongebob:",
		AutoComplete:     true,
		AutoCompleteDesc: "tYpE wHaT yOu WaNt To sAy In MocKInG tOnE :spongebob:",
	})
	return err
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	output := fmt.Sprintf("%s :spongebob:", dOiT(strings.TrimPrefix(args.Command, "/spongebob ")))
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
		Text:         output,
	}, nil
}

func mockingText(message string) (string, bool) {
	message, next := mockingTextOnce(message)
	if next == -1 {
		return message, false
	}

	cursor := next
	for cursor < len(message) {
		nextMessage, next := mockingTextOnce(message[cursor:])
		if next == -1 {
			return message, true
		}
		message = message[:cursor] + nextMessage
		cursor += next
	}

	return message, true
}

func dOiT(input string) string {
	base := rand.Int() % 2
	contentSlice := strings.Split(input, "")
	for i := range contentSlice {
		if rand.Int()%10 > 9 {
			continue
		}
		if i%2 == base {
			contentSlice[i] = strings.ToUpper(contentSlice[i])
		} else {
			contentSlice[i] = strings.ToLower(contentSlice[i])
		}
	}
	return strings.Join(contentSlice, "")
}

func mockingTextOnce(message string) (string, int) {
	matchIndex := strings.Index(message, spongeBob)
	if matchIndex == -1 {
		return message, -1
	}
	prelude := message[:matchIndex]
	postlude := message[matchIndex:]
	closingMatchIndex := strings.Index(postlude[spongeBobLen:], spongeBob)
	if closingMatchIndex == -1 {
		return message, -1
	}
	contentSlice := strings.Split(postlude[spongeBobLen:closingMatchIndex+spongeBobLen], "")
	postlude = postlude[spongeBobLen*2+closingMatchIndex:]
	base := rand.Int() % 2
	for i := range contentSlice {
		if rand.Int()%10 > 9 {
			continue
		}
		if i%2 == base {
			contentSlice[i] = strings.ToUpper(contentSlice[i])
		} else {
			contentSlice[i] = strings.ToLower(contentSlice[i])
		}
	}
	message = strings.Join([]string{
		prelude,
		spongeBob,
		strings.Join(contentSlice, ""),
		spongeBob,
		postlude,
	}, "")
	return message, len(prelude) + len(contentSlice) + spongeBobLen*2
}

func (p *Plugin) spongeBobMeBoy(c *plugin.Context, post *model.Post) (*model.Post, string) {
	message, changed := mockingText(post.Message)
	if changed {
		post.Message = message
		return post, ""
	}

	return nil, ""
}

func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	return p.spongeBobMeBoy(c, post)
}

func (p *Plugin) MessageWillBeUpdated(c *plugin.Context, post *model.Post) (*model.Post, string) {
	return p.spongeBobMeBoy(c, post)
}
