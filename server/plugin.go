// See https://developers.mattermost.com/extend/plugins/server/reference/
package main

import (
	"fmt"
	"strings"
	"sync"

	"math/rand"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

const spongeBob string = ":spongebob:"
const spongeBobLen int = len(spongeBob)

func mockingText(message string) (string, bool) {
	message, next := mockingTextOnce(message)
	if next == -1 {
		return message, false
	}

	// cursor := next
	// for cursor < len(message) {
	// 	nextMessage, next := mockingTextOnce(message[cursor:])
	// 	if next == -1 {
	// 		return message, true
	// 	}
	// 	cursor += next
	// 	message = message[:next] + nextMessage
	// }

	return message[:next] + "][" + message[next:], true
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
	for i := range contentSlice {
		if rand.Int()%2 == 0 {
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
	fmt.Println(message)
	fmt.Println(closingMatchIndex)
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
