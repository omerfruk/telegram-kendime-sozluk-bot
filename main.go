package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var token = "5230692118:AAEsx_ATWFmzwRFyHtM7ErA_UPjLu9ld600"

type Sozluk struct {
	MetinIng   string
	OrnekCumle string
}

func main() {

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Token)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	go TimeMachine(bot)

	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help", "h":
				msg.Text = "/introduceYourself : Introduce yourself\n" +
					"/newWord : Sends you a new word\n" +
					"Reply message with \"/explain\" (Very soon)"
			case "introduceYourself":
				msg.Text = "Hi my name is Kendine Sözlük bot :)\n" +
					"My goal is to teach you new words\n" +
					"I'll send you a word every hour between 10:00 and 22:00\n" +
					"If you want a new word without waiting, what you should do is; /newWord\n" +
					"My author is @OmerFrukTasdemir"
			case "newWord":
				sozcuk := SozcukGetirici()
				msg.Text = fmt.Sprintf("%s\n\n%s", sozcuk.MetinIng, sozcuk.OrnekCumle)
			default:
				msg.Text = "I don't know that command"
			}

			if _, err = bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}

func TimeMachine(bot *tgbotapi.BotAPI) {
	ticker := time.NewTicker(1 * time.Minute)
	for _ = range ticker.C {
		if time.Now().Hour() < 22 && time.Now().Hour() > 10 {
			sozcuk := SozcukGetirici()
			msg := tgbotapi.NewMessage(-702366902, fmt.Sprintf("%s\n\n%s",
				sozcuk.MetinIng,
				sozcuk.OrnekCumle,
			))
			bot.Send(msg)
			time.Sleep(1 * time.Hour)
		}
	}
}

func SozcukGetirici() Sozluk {
	byteSozluk, err := ioutil.ReadFile("sozluk.json")
	if err != nil {
		log.Println("Sozcuk dosyadan okunurken hata aldı")
	}
	var sozlukler []Sozluk
	err = json.Unmarshal(byteSozluk, &sozlukler)
	if err != nil {
		log.Println("sozcukler unmarshal edilirken hata olustu ")
	}
	rnd := rand.Intn(len(sozlukler)-0) + 0
	sozcuk := sozlukler[rnd]

	return sozcuk
}

func kelimeGet() {
	//teknik kelimler için link
	//https://www.fullofenglish.com/mesleki-ingilizce-bilisim/ingilizce-yazilim-terimleri/
	res, err := http.Get("https://gonaturalenglish.com/1000-most-common-words-in-the-english-language/")
	if err != nil {
		fmt.Println(err)
	}
	if res.StatusCode != 200 {
		fmt.Println("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	var sozluk []Sozluk
	defer res.Body.Close()
	doc.Find("ol > li").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		if i >= 3 {
			title = strings.ReplaceAll(title, "\"", "")
			title = strings.ReplaceAll(title, "“", "")
			title = strings.ReplaceAll(title, "”", "")
			cumleler := strings.Split(title, "–")
			if len(cumleler) == 2 {
				cumleler[0] = strings.TrimSpace(cumleler[0])
				cumleler[1] = strings.TrimLeft(cumleler[1], " ")
				cumleler[1] = strings.Replace(cumleler[1], cumleler[0], fmt.Sprintf("%s", cumleler[0]), -1)
				sozluk = append(sozluk, Sozluk{
					MetinIng:   cumleler[0],
					OrnekCumle: cumleler[1],
				})
			}
		}

	})

	sozlukJson, _ := json.MarshalIndent(&sozluk, "", " ")

	ioutil.WriteFile("sozluk.json", sozlukJson, 0644)

}
