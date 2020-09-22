package main

import (
	"bytes"
	"ding/ding"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Global struct {
		RoutePath  string `yaml:"routePath"`
		ListenPort string `yaml:"listenPort"`
	} `yaml:"global"`

	Ding struct {
		Token     string `yaml:"token"`
		Title     string `yaml:"title"`
		MsgFormat struct {
			MsgType    string `yaml:"msgType"`
			MsgHead    string `yaml:"msgHead"`
			MsgBodyTpl string `yaml:"msgBodyTpl"`
		} `yaml:"msgFormat"`
	} `yaml:"ding"`
}

type Notification struct {
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	Alerts            []Alert           `json:"alerts"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`

	Version  string `json:"version"`
	GroupKey string `json:"groupKey"`
}

type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

func loadConfig() *Config {
	cfg := new(Config)
	fs, err := os.Open("./ding.yml")
	defer fs.Close()

	if err != nil {
		log.Panicln(err)
	}

	cfgBuf, err := ioutil.ReadAll(fs)
	if err != nil {
		log.Panicln(err)
	}
	yaml.Unmarshal(cfgBuf, cfg)

	return cfg
}

func resiveNotif(c *gin.Context) (*Notification, error) {
	notif := new(Notification)
	if err := c.BindJSON(&notif); err != nil {
		log.Println("error " + err.Error())
		c.JSON(http.StatusOK, gin.H{"message": "error " + err.Error()})
		return nil, err
	}
	return notif, nil
}

func formatMsg(cfg *Config, notif *Notification) string {
	msg := cfg.Ding.MsgFormat.MsgHead
	tpl, _ := template.ParseFiles(cfg.Ding.MsgFormat.MsgBodyTpl)
	buf := new(bytes.Buffer)
	for _, alert := range notif.Alerts {
		if err := tpl.Execute(buf, struct {
			Alert Alert
		}{alert}); err != nil {
			log.Println(err)
		}

		msg += buf.String()
		buf.Reset()
	}
	return msg
}

func main() {
	cfg := loadConfig()
	bot := ding.NewRobot(cfg.Ding.Token)
	router := gin.Default()
	router.POST(cfg.Global.RoutePath, func(c *gin.Context) {
		notif, err := resiveNotif(c)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"message": "receive alert notification message failed!"})
			log.Println(err)
		}

		msg := formatMsg(cfg, notif)

		c.JSON(http.StatusOK, gin.H{"message": "receive alert notification message successful!"})
		bot.Say(ding.Msg{
			MsgType: cfg.Ding.MsgFormat.MsgType,
			Markdown: ding.Markdown{
				Title: cfg.Ding.Title,
				Text:  msg,
			}, Text: ding.Text{
				Content: msg,
			},
		})
	})

	if err := router.Run(":" + cfg.Global.ListenPort); err != nil {
		log.Println(err)
	}
}
