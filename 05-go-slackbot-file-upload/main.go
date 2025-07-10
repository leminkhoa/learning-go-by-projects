package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	channelID := os.Getenv("CHANNEL_ID")
	fileArr := []string{"ZIPL.pdf"}
	for i := 0; i < len(fileArr); i++ {

		fileInfo, err := os.Stat(fileArr[i])
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		params := slack.UploadFileV2Parameters{
			Channel:  channelID,
			File:     fileArr[i],
			Filename: string(fileInfo.Name()),
			FileSize: int(fileInfo.Size()),
		}
		file, err := api.UploadFileV2(params)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("ID: %s, title: %s\n", file.ID, file.Title)
	}
}
