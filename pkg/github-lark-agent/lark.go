package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	"log"
	"net/http"
	"strings"
)

func sendMessageToLark(apiUrl string, title string, content string) error {

	body := `{
	"msg_type": "interactive",
	"card": {
		"config": {
			"wide_screen_mode": true
		},
		"header": {
			"title": {
				"tag": "plain_text",
				"content": ""
			},
			"template": "blue"
		},
		"elements": [{
			"tag": "markdown",
			"content": null
		}]
	}
}`
	body, err := sjson.Set(body, "card.header.title.content", title)
	if err != nil {
		return err
	}
	body, err = sjson.Set(body, "card.elements.0.content", content)
	if err != nil {
		return err
	}

	resp, err := http.Post(apiUrl, "application/json", strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
		return err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}

	bodyString := string(bodyBytes)
	log.Println("lark resp:", resp.StatusCode, "body: ", bodyString)

	status := gjson.Get(bodyString, "StatusMessage").Int()

	if status == 0 {
		return nil
	}

	fmt.Printf("post failed, err code :%n\n", status)
	return err

}
