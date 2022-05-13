package main

import (
	"fmt"
	"github.com/go-playground/webhooks/v6/github"
)

func makeGithubPushMessage(body github.PushPayload) string {

	ref := body.Ref
	pusher := body.Pusher.Name
	repo := body.Repository.FullName

	branch := refToBranch(ref)

	text := ""
	text += "Repo: " + repo + "\n"
	text += "Branch: " + branch + "\n"
	text += "Pusher: " + pusher + "\n"
	text += "Commits:\n"

	commits := body.Commits

	if len(commits) > 0 {
		for _, commit := range commits {

			author := commit.Author
			committer := commit.Committer.Name
			message := commit.Message
			id := commit.ID
			url := commit.URL

			text += fmt.Sprintf("  %s author: %s, committer: %s :%s [查看改变](%s)\n",
				id, author, committer, message, url,
			)
		}
	}

	return text

}

func refToBranch(refs string) string {
	var lastIndex int
	for i, cha := range refs {
		if cha == '/' {
			lastIndex = i
		}
	}

	rs := []rune(refs)

	return string(rs[lastIndex+1:])
}
