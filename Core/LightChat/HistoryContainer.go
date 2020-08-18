package LightChat

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
)

type HistoryContainer struct {
	SiteUrl string

	lastMessageId string
}

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