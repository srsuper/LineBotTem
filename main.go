// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"html/template"
	"strings"
	"vendor/linebot"
	"vendor/golang.org"
	"vendor/examples"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main(){
	var err error
	bot, err = linebot.New(os.Getenv("3eda0f01f6b2fe05aace973e39707edc"), os.Getenv("XuAPgE5eH13Hbgj7mSSCmqe5wheTgVDhiE805ypPKx1hyHXCXLgshl02rpLCe+rUUVTfBE6SkoXrkRD0c1omm6o8RFZMgCETtwF7nDTKSg3PDQG6OIHE2npC1e3YfWXhvBMcXBwFrF5zE8s9T83cgQdB04t89/1O/w1cDnyilFU="))
	log.Println("Bot:", bot, " err:", err)
	//Default: Hello World!
	http.HandleFunc("/", HelloServer)
	//User input form
	http.HandleFunc("/InputMsg", InputMsgHandler)
	//Send message by user input
	http.HandleFunc("/SendMsg", SendMsgHandler)
	//1. Can Reply message when use Line
	//2. Use html Get to Send message
	//EXAMPLE: https://loaclhost/SendMsg?type=SendMsg&uid=[UserID]&msg=[Message you want to send]
	http.HandleFunc("/SendLineMsg", SendLineMsgHandler)

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}
func SendLineMsgHandler(w http.ResponseWriter, r *http.Request) {
	//Parse HttpRequest
	events, err := bot.ParseRequest(r)
	log.Println("Events:", events, " err:", err)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
			fmt.Println(w, "\n400\n")
			fmt.Fprintln(w, "\n400\n")
		}else {
			w.WriteHeader(500)
			fmt.Println(w, "\n500\n")
			fmt.Fprintln(w, "\n500\n")
		}
		return
	}else {
		for _, event := range events 	{
			if event.Type == linebot.EventTypeMessage {
				log.Println("Your Line User ID: \n  " + event.Source.UserID + "\n")
				log.Println("Your Line Group ID: \n  " + event.Source.GroupID + "\n")
				userid := "Your Line User ID: \n  " + event.Source.UserID + "\n"
				groupid := "Your Line Group ID: \n  " + event.Source.GroupID + "\n"
				switch message := event.Message.(type) {
					case *linebot.TextMessage:
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(userid + groupid + "Message: \n  "+ message.Text)).Do(); err != nil {
				            log.Print(err)
						}
				}
			}
		}
	}
}
func SendMsgHandler(w http.ResponseWriter, r *http.Request) {
	//Parse HttpRequest
	//notice: if not use ParseForm, you can't get the values
	r.ParseForm()
	val := [3][2]string{{"type", ""}, {"uid", ""}, {"msg", ""}}

	//Put the values to array
	i := 0
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
		if k == "type" {
			val[0][1] = strings.Join(v, "")
		}
		if k == "uid" {
			val[1][1] = strings.Join(v, "")
		}
		if k == "msg" {
			val[2][1] = strings.Join(v, "")
		}
		i++
	}
	//Print values on the output Page
	for i := 0; i < 3; i++ {
		fmt.Fprint(w, val[i][0])
		fmt.Fprint(w, " :\n")
		fmt.Fprintln(w, val[i][1])
		fmt.Fprintln(w, "")
	}
	//Send Message to user
	if val[0][1] == "SendMsg" {
		bot, err := linebot.New(os.Getenv("3eda0f01f6b2fe05aace973e39707edc"), os.Getenv("XuAPgE5eH13Hbgj7mSSCmqe5wheTgVDhiE805ypPKx1hyHXCXLgshl02rpLCe+rUUVTfBE6SkoXrkRD0c1omm6o8RFZMgCETtwF7nDTKSg3PDQG6OIHE2npC1e3YfWXhvBMcXBwFrF5zE8s9T83cgQdB04t89/1O/w1cDnyilFU="))
		if err != nil {
		}
		if _, err := bot.PushMessage(val[1][1], linebot.NewTextMessage(val[2][1])).Do(); err != nil {
		}
	}
}
func InputMsgHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
    if r.Method == "GET" {
        t, _ := template.ParseFiles("InputLineMsg.html")
        t.Execute(w, nil)
		//v := url.Values{}
    } else {
        r.ParseForm()
        fmt.Println("User ID:", r.Form["method"])
        fmt.Println("User ID:", r.Form["uid"])
        fmt.Println("Message:", r.Form["msg"])
    }
}
func HelloServer(w http.ResponseWriter, r *http.Request) {
	//Parse the parameter from httpReuest
	//notice: if not use ParseForm, you can't get the values
	r.ParseForm()
	//Print on the Server
    fmt.Println(r.Form)
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
	fmt.Fprintln(w, "hello, world!\n")
}
