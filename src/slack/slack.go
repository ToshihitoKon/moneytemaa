package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ToshihitoKon/moneytemaa/src/constants"
	mydb "github.com/ToshihitoKon/moneytemaa/src/db"
	"github.com/gorilla/mux"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var (
	api = slack.New(constants.SlackBotToken)
)

func SetHandler(router *mux.Router) {
	router.HandleFunc("/slack/event", HandlerSlackEvent)
}

func HandlerSlackEvent(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sv, err := slack.NewSecretsVerifier(r.Header, constants.SlackSigningSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := sv.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := sv.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	case slackevents.CallbackEvent:
		innerEventAction(eventsAPIEvent.InnerEvent)
	}
}

func innerEventAction(innerEvent slackevents.EventsAPIInnerEvent) {
	switch ev := innerEvent.Data.(type) {
	case *slackevents.ReactionAddedEvent:
		reactionAddedEvent, ok := innerEvent.Data.(*slackevents.ReactionAddedEvent)
		if !ok {
			log.Println("err: slackevents.ReactionAddedEvent")
			return
		}
		log.Println("ReactionAddedEvent: ", reactionAddedEvent)

	case *slackevents.MessageEvent:
		messageEvent, ok := innerEvent.Data.(*slackevents.MessageEvent)
		if !ok {
			log.Println("err: slackevents.MessageEvent")
			return
		}
		if messageEvent.User == constants.SlackBotUserId {
			return
		}
		messageEventAction(messageEvent)

	case *slackevents.AppMentionEvent:
		candidate := strings.Split(Dummytext, "\n")
		rand.Seed(time.Now().UnixNano())
		ret := candidate[rand.Intn(len(candidate))]
		_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText(ret, false))
		if err != nil {
			log.Println("slackevents.AppMentionEvent: ", err.Error())
		}

	default:
		log.Println("default: ", innerEvent.Type)
	}
}

func messageEventAction(messageEvent *slackevents.MessageEvent) {
	var err error
	if strings.HasPrefix(messageEvent.Text, "散財") {
		// 分割
		separate := regexp.MustCompile(`[\n| ]+`)
		text := separate.Split(messageEvent.Text, -1)
		if len(text) < 3 {
			_, _, err = api.PostMessage(messageEvent.Channel, slack.MsgOptionText("[金額] [メモ]", false))
			if err != nil {
				log.Println("messageEventAction: ", err.Error())
			}
			return
		}

		// 金額
		price, err := strconv.Atoi(text[1])
		if err != nil {
			log.Println("err: strconv.Atoi: ", err.Error())
		}

		// メモ
		comment := text[2]

		slackResponse := fmt.Sprint(
			"メモった",
			"\n額: ", strconv.Itoa(price),
			"\nメモ: ", comment,
		)

		err = mydb.SlackInsertTransaction(price, comment, messageEvent.User, messageEvent.Channel, messageEvent.TimeStamp)
		if err != nil {
			slackResponse = fmt.Sprint("失敗したわ ", err.Error())
		}

		_, _, err = api.PostMessage(messageEvent.Channel, slack.MsgOptionText(slackResponse, false))
		if err != nil {
			log.Println("messageEventAction: ", err.Error())
		}
	}
}
