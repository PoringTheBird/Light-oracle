package LightChat

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"time"
)

const historyLoadInterval = 10

type MessageLoadHandler interface {
	OnNewMessagesLoaded(messages []Message)
}

type HistoryContainer struct {
	SiteUrl string
	LoadHandler MessageLoadHandler

	lastMessageId string

	historyObserving chan struct{}
}

// Actions

func (container *HistoryContainer) LoadHistory() (*[]Message, error) {
	url := fmt.Sprintf("%s?mode=getshouts&jal_lastId=%s", container.SiteUrl, container.lastMessageId)

	resp, err := http.Get(url)
	if err != nil { return nil, err }

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil { return nil, err }

	regex, err := regexp.Compile("([\\d]+)\\|\\|([\\d]+)###([^#]+)###[^#]+###([^#]+)###[^#]+###[^#]+###")
	if err != nil { return nil, err}

	submatches := regex.FindAllStringSubmatch(string(content), -1)

	var messages = new([]Message)

	for _, match := range submatches{
		msg := Message{MessageId: match[1], SenderId: match[2], SenderName: match[3], Text: match[4]}

		*messages = append(*messages, msg)
	}

	sort.Slice(*messages, func(i, j int) bool {
		arr := *messages
		return arr[i].MessageId < arr[j].MessageId
	})

	return messages, nil
}

func (container *HistoryContainer) StartChatHistoryObserving(interval int) {
	historyReload := time.NewTicker(time.Duration(interval) * time.Millisecond)
	container.historyObserving = make(chan struct{})

	go func() {
		for {
			select {
			case <- historyReload.C:
				err := container.sendNewMessages()

				if err != nil {
					log.Println("History reload stopped. Reason: ", err)
					historyReload.Stop()
					return
				}
			case <- container.historyObserving:
				historyReload.Stop()
				return
			}
		}
	}()
}

func (container *HistoryContainer) StopChatHistoryObserving() {
	close(container.historyObserving)
}

// Internal actions

func (container *HistoryContainer) sendNewMessages() error {
	messages, err := container.LoadHistory()
	if err != nil { return err }

	var newMessages []Message

	for _, msg := range *messages {
		if msg.MessageId > container.lastMessageId {
			newMessages = append(newMessages, msg)
		}
	}

	messagesCount := len(newMessages)

	if messagesCount > 0 {
		if len(container.lastMessageId) > 0 {
			container.LoadHandler.OnNewMessagesLoaded(newMessages)
		}

		container.lastMessageId = newMessages[messagesCount - 1].MessageId
	}

	return nil
}