package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nlopes/slack"
)

/*
   TODO: Change @BOT_NAME to the same thing you entered when creating your Slack application.
   NOTE: command_arg_1 and command_arg_2 represent optional parameteras that you define
   in the Slack API UI
*/
const helpMessage = "type in '@goemoji what happens if robots take over the world?'"

var apiKey = os.Getenv("EMOJI_API_KEY")

// EmojiResult contains API data
type EmojiResult struct {
	Result  bool   `json:"result"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Text    string `json:"text"`
}

/*
   CreateSlackClient sets up the slack RTM (real-timemessaging) client library,
   initiating the socket connection and returning the client.
   DO NOT EDIT THIS FUNCTION. This is a fully complete implementation.
*/
func CreateSlackClient(apiKey string) *slack.RTM {
	api := slack.New(apiKey)
	rtm := api.NewRTM()
	go rtm.ManageConnection() // goroutine!
	return rtm
}

/*
   RespondToEvents waits for messages on the Slack client's incomingEvents channel,
   and sends a response when it detects the bot has been tagged in a message with @<botTag>.!
*/
func RespondToEvents(slackClient *slack.RTM) {
	for msg := range slackClient.IncomingEvents {
		fmt.Println("Event Received: ", msg.Type)
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			botTagString := fmt.Sprintf("<@%s> ", slackClient.GetInfo().User.ID)
			if !strings.Contains(ev.Msg.Text, botTagString) {
				continue
			}
			message := strings.Replace(ev.Msg.Text, botTagString, "", -1)
			sendEmoji(slackClient, message, ev.Channel)
			sendHelp(slackClient, message, ev.Channel)
		default:

		}
	}
}

// sendHelp is a working help message, for reference.
func sendHelp(slackClient *slack.RTM, message, slackChannel string) {
	if strings.ToLower(message) != "help" {
		return
	}
	slackClient.SendMessage(slackClient.NewOutgoingMessage(helpMessage, slackChannel))
}

func sendEmoji(slackClient *slack.RTM, message, slackChannel string) {
	cmd := strings.ToLower(message)
	url := "https://api.ritekit.com/v1/emoji/auto-emojify?text=" + cmd + "&client_id=" + apiKey

	fmt.Printf("%v", url)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	emojifiedText, _ := ioutil.ReadAll(resp.Body)
	emojiResult := EmojiResult{}
	err = json.Unmarshal(emojifiedText, &emojiResult)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", emojiResult)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.Status)
	fmt.Println(string(respBody))

	println(emojiResult.Message)
	println("[RECV] sendResponse:", cmd)
	println("[API] emojifiedText:", emojiResult.Message, emojiResult.Text)
	slackClient.SendMessage(slackClient.NewOutgoingMessage(emojiResult.Text, slackChannel))

}
