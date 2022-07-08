package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/enescakir/emoji"
	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

func displayCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("command events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	//load env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	slackBotToken, exists := os.LookupEnv("SLACK_BOT_OAUTH-TOKEN")
	slackAppToken, isAvail := os.LookupEnv("SLACK-USER-OAUTH-TOKEN")
	//check if slackBottoken is available
	if !exists {
		log.Fatal("no env variable found for the SLACK_BOT_OAUTH-TOKEN")
	}
	//check if slackBottoken is available
	if !isAvail {
		log.Fatal("no env variable found for the SLACK-USER-OAUTH-TOKEN")
	}
	//create new bot client
	bot := slacker.NewClient(slackBotToken, slackAppToken)

	go displayCommandEvents(bot.CommandEvents())

	//create new bot command
	bot.Command("year of birth is <year>", &slacker.CommandDefinition{
		Description: "Age calculator",
		Example:     "my year of birth is 2022",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			//get the param
			yearOfBirthInString := request.Param("year")
			//convert the param to strin
			yearOfBirthInInt, strconvError := strconv.Atoi(yearOfBirthInString)
			var r string
			if strconvError != nil {
				println("error")
			}
			//get the current year
			date := time.Now().Year()
			age := date - yearOfBirthInInt
			if age < 0 {
				r = fmt.Sprintf("You are from the future %v", emoji.RollingOnTheFloorLaughing)
				response.Reply(r)
				return
			}
			r = fmt.Sprintf("age is %d", age)
			//response to the event
			response.Reply(r)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//Listen receives events from Slack and each is handled as needed
	AnotherErr := bot.Listen(ctx)

	if AnotherErr != nil {
		log.Fatal(err)
	}

}
