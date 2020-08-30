package Discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"main/Core/Discord/Entities"
	"net/http"
	"sort"
)

type ApiClient struct {
	BaseUrl string
	DiscordToken string
}

// Actions

func (api ApiClient) LoadChannel(channelId string) (*Entities.Channel, error) {
	channel := new(Entities.Channel)

	path := fmt.Sprintf("/channels/%s", channelId)
	err := api.makeRequest("GET", path, nil, channel)

	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (api ApiClient) LoadMessages(channelId string) (*[]Entities.IncomingMessage, error) {
	messages := new([]Entities.IncomingMessage)

	path := fmt.Sprintf("/channels/%s/messages", channelId)
	err := api.makeRequest("GET", path, nil, messages)

	if err != nil {
		return nil, err
	}

	// Discord API returns unsorted messages list
	sort.Slice((*messages), func(i, j int) bool {
		return (*messages)[j].Timestamp.Sub((*messages)[i].Timestamp) >= 0
	})

	return messages, nil
}

func (api ApiClient) SendMessage(text string, channelId string) (*Entities.IncomingMessage, error) {
	outcomingMessage := Entities.OutcomingMessage{Content: text, Tts: false, Embed: nil}
	sentMessage := new (Entities.IncomingMessage)

	path := fmt.Sprintf("/channels/%s/messages", channelId)
	err := api.makeRequest("POST", path, outcomingMessage, sentMessage)

	if err != nil {
		return nil, err
	}

	return sentMessage, nil
}

// Api Core

func (api ApiClient) makeRequest(
	method string,
	path string,
	params interface{},
	resultModel interface{}) error {

	client := &http.Client{}

	url := api.BaseUrl + path

	var buffer io.Reader

	if params != nil {
		data, err := json.Marshal(params)

		if err != nil {
			return err
		}

		buffer = bytes.NewBuffer(data)
	} else {
		buffer = nil
	}

	req, err := http.NewRequest(method, url, buffer)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bot " + api.DiscordToken)

	resp, err := client.Do(req)

	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		text, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(text))
	}

	decodingErr := json.NewDecoder(resp.Body).Decode(resultModel)

	if decodingErr != nil {
		return decodingErr
	}

	return err
}