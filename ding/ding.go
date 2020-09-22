package ding

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const URL = "https://oapi.dingtalk.com/robot/send?access_token="

type Robot struct {
	token string
}

type Msg struct {
	MsgType string `json:"msgtype"`
	Text    Text   `json:"text"`

	At       At       `json:"at"`
	Link     Link     `json:"link"`
	Markdown Markdown `json:"markdown"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type Text struct {
	Content string `json:"content"`
}

type Link struct {
	Text      string `json:"text"`
	Title     string `json:"title"`
	PicUrl    string `json:"picUrl"`
	MssageUrl string `json:"messageUrl"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type ErrorMsg struct {
	Errcode string `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func NewRobot(token string) *Robot {
	r := new(Robot)
	r.token = token
	return r
}

func (r *Robot) Say(t Msg) error {
	//编码 json
	data, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		return err
	}

	//请求 钉钉
	buf := bytes.Buffer{}
	buf.Write(data)
	resp, err := http.Post(URL+r.token, "application/json", &buf)
	if err != nil {
		log.Println(err)
		return err
	}

	//解析 钉钉返回错误信息
	errMsg := new(ErrorMsg)
	buf.ReadFrom(resp.Body)
	json.Unmarshal(buf.Bytes(), errMsg)
	if errMsg.Errcode != "0" {
		log.Println(errMsg.Errmsg)
		return errors.New(errMsg.Errmsg)
	}

	return nil
}
