package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type SendData struct {
	Subject string `json:"subject"`
	// Body    BodyContent `json:"body"`
	Start StartDate `json:"start"`
	End   EndDate   `json:"end"`
	// Attendees []AttendeesData `json:"attendees"`
}

type BodyContent struct {
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
}

type StartDate struct {
	DateTime string `json:"dateTime"`
	TimeZone string `json:"timeZone"`
}

type EndDate struct {
	DateTime string `json:"dateTime"`
	TimeZone string `json:"timeZone"`
}

type AttendeesData struct {
	EmailAddress EmailAddress `json:"emailAddress"`
	AType        string       `json:"type"`
}

type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

func main() {
	token := ""
	send(token)
}

func send(token string) {
	//content := BodyContent{
	//	ContentType: "HTML",
	//	Content:     "請假事由～ TEST",
	//}
	sdate := StartDate{
		DateTime: time.Now().UTC().Format("2006-01-02T03:04:05"),
		TimeZone: "UTC",
	}
	edate := EndDate{
		DateTime: time.Now().UTC().Format("2006-01-02T03:04:05"),
		TimeZone: "UTC",
	}

	var attDate []AttendeesData
	tmp := AttendeesData{
		EmailAddress: EmailAddress{
			Address: "",
			Name:    "",
		},
		AType: "required",
	}

	attDate = append(attDate, tmp)

	var postData = SendData{
		Subject: "test on leave",
		// Body:    content,
		Start: sdate,
		End:   edate,
		// Attendees: attDate,
	}
	body, _ := json.Marshal(postData)

	fmt.Println(string(body))
	req, _ := http.NewRequest("POST", "https://graph.microsoft.com/v1.0/me/events", strings.NewReader(string(body)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Println("save topic failed", err.Error())
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		jsonStr := string(body)
		fmt.Println("Response: ", jsonStr)

	} else {
		fmt.Println("Get failed with error: ", resp.Status)
	}
}
