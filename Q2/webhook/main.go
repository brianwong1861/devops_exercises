package main

import (
	"os"
	"time"
	"log"
	"fmt"
	"net/http"
	_ "reflect"
	"gopkg.in/go-playground/webhooks.v5/github"
	"github.com/nlopes/slack"
	_ "github.com/joho/godotenv"
)
// init() is uncommented if you run it in Docker
// func init(){
// 	_ = godotenv.Load(".env")
// }


func epochToHeumanReadable() time.Time {
	currentTime := time.Now().Unix()
	return time.Unix(currentTime, 0)
}

func SendToSlack(commit, msg string){
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Pretext: "Repository Notification",
		Text: "Commit ID:\t" + commit + "\n" + "Commit Messages:\t" + msg,
	}
	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := api.PostMessage(os.Getenv("CHANNEL_ID"),"Push Notification", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
}


func main() {
	hook, _ := github.New(github.Options.Secret(os.Getenv("SECRET")))
	
	http.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.ReleaseEvent, github.PullRequestEvent, github.PushEvent, github.PullRequestReviewEvent, github.PullRequestReviewCommentEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
			}
		}
		switch payload.(type) {
			
		case github.PushPayload:
			pushPayload := payload.(github.PushPayload)
			SendToSlack(pushPayload.HeadCommit.ID, pushPayload.HeadCommit.Message)
			log.Println(pushPayload.HeadCommit.ID + "===>" + pushPayload.HeadCommit.Message)


		}
	})
	fmt.Printf("Serving on %s\n", os.Getenv("PORT"))
	http.ListenAndServe(":"+ string(os.Getenv("PORT")), nil)
}
