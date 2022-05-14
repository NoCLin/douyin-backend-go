package main

import (
	"fmt"
	"net/http"
	"os"
)

import (
	"github.com/go-playground/webhooks/v6/github"
)

func main() {

	SecretKey := os.Getenv("AGENT_SECRET")
	LarkApiUrl := os.Getenv("LARK_API_URL")

	fmt.Println("SK", SecretKey)
	fmt.Println("LarkApiUrl", LarkApiUrl)

	hook, _ := github.New(github.Options.Secret(SecretKey))

	http.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PushEvent, github.ReleaseEvent, github.PullRequestEvent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		switch payload.(type) {

		case github.PushPayload:
			push := payload.(github.PushPayload)

			title := "新的 Github 推送"
			content := makeGithubPushMessage(push)
			fmt.Println(content)

			err = sendMessageToLark(LarkApiUrl, title, content)

			_, _ = fmt.Fprintf(w, "OK")
			return

		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", pullRequest)
			return
		}
	})

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("listen error")
		return
	}
}
