package main

import (
	"bytes"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Variable struct {
	KeyOrId string `json:"keyOrId"`
}
type GraphQlQuery struct {
	Query     string   `json:"query"`
	Variables Variable `json:"variables"`
}

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func getJiraItem(jiraId string, token string) string {
	jiraGraphQlUri := "https://development-tools.global.dish.com/services/jira/api/v1/graphql"
	vars := Variable{KeyOrId: jiraId}
	query := GraphQlQuery{
		Query:     "query($keyOrId:String!){ issueByKeyOrId(keyOrId:$keyOrId){ url fields }}",
		Variables: vars,
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		log.Println("Error formating query:", err)
		return ""
	}

	req, err := http.NewRequest("POST", jiraGraphQlUri, bytes.NewBuffer(queryBytes))
	if err != nil {
		log.Println("Error building request:", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")
	bearer := "Bearer " + token
	req.Header.Add("Authorization", bearer)
	jiraClient := &http.Client{}
	jiraResponse, err := jiraClient.Do(req)
	if err != nil {
		log.Println("Error when getting jira item:", err)
		return ""
	}
	defer jiraResponse.Body.Close()

	log.Println("Response Status:", jiraResponse.Status)
	if jiraResponse.StatusCode == http.StatusOK {
		jiraItem, _ := io.ReadAll(jiraResponse.Body)
		return string(jiraItem)
	}
	return ""
}

func main() {
	a := app.New()
	w := a.NewWindow("JiraWidgetLite")
	w.Resize(fyne.NewSize(400, 400))

	summaryText := widget.NewLabel("")
	startText := widget.NewLabel("")
	stopText := widget.NewLabel("")
	durationText := widget.NewLabel("")

	startTime := time.Now()
	stopTime := time.Now()
	duration := time.Duration(0)

	startButton := widget.NewButton("Start", func() {
		startTime = time.Now()
		startText.SetText(startTime.Format(time.RFC3339))
		log.Println(startTime.Format(time.RFC3339))
	})
	stopButton := widget.NewButton("Stop", func() {
		stopTime = time.Now()
		stopText.SetText(stopTime.Format(time.RFC3339))
		log.Println(stopTime.Format(time.RFC3339))
		duration = stopTime.Sub(startTime)
		durationText.SetText(duration.String())
	})

	resetButton := widget.NewButton("Reset", func() {

		startText.SetText("")
		stopText.SetText("")
		durationText.SetText("")
	})

	buttonContainer := container.NewHBox()

	buttonContainer.Add(startButton)
	buttonContainer.Add(stopButton)
	buttonContainer.Add(resetButton)

	tokenContainer := container.NewVBox()
	token := widget.NewMultiLineEntry()

	jiraId := widget.NewEntry()
	jiraButton := widget.NewButton("pull", func() {
		log.Println("id:", jiraId.Text)
		summaryText.SetText("")
		jiraItem := getJiraItem(jiraId.Text, token.Text)
		bytJiraItem := []byte(jiraItem)
		var mapData map[string]interface{}
		json.Unmarshal(bytJiraItem, &mapData)
		data := mapData["data"].(map[string]interface{})
		issue := data["issueByKeyOrId"].(map[string]interface{})
		fields := issue["fields"].(map[string]interface{})
		summary := fields["Summary"].(string)
		log.Println(summary)
		summaryText.SetText(summary)

	})

	jiraIdContainer := container.NewWithoutLayout(jiraId, jiraButton)
	jiraId.Resize(fyne.NewSize(200, 32))
	jiraId.Move(fyne.NewPos(0, 0))
	jiraId.Wrapping = fyne.TextWrapOff
	jiraId.Scroll = container.ScrollNone
	jiraButton.Resize(fyne.NewSize(40, 32))
	jiraButton.Move(fyne.NewPos(208, 0))

	tokenUrl, _ := url.Parse("https://development-tools.global.dish.com/api/account/access-token?name=openid")

	tokenContainer.Add(container.NewHBox(widget.NewLabel("Gateway Token:"), widget.NewHyperlink("MyToken", tokenUrl)))

	tokenContainer.Add(container.NewVBox(token))
	tokenContainer.Add(widget.NewLabel("Jira Id"))
	tokenContainer.Add(jiraIdContainer)

	jiraItemContainer := container.NewVBox()

	jiraItemContainer.Add(container.NewHBox(widget.NewLabel("Summary"), summaryText))
	jiraItemContainer.Add(container.NewHBox(widget.NewLabel("Start"), startText))
	jiraItemContainer.Add(container.NewHBox(widget.NewLabel("End"), stopText))
	jiraItemContainer.Add(container.NewHBox(widget.NewLabel("Duration"), durationText))

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "", Widget: tokenContainer},
			{Text: "", Widget: buttonContainer},
			{Text: "", Widget: jiraItemContainer},
		},
		OnSubmit: func() { // update to push worktime to item in jira
			log.Println("token:", token.Text)
			log.Println("id:", jiraId.Text)
			log.Println("start:", startTime.Format(time.RFC3339))
			log.Println("stop:", stopTime.Format(time.RFC3339))
			log.Println("duration:", duration)

		},
	}
	w.SetContent(form)
	w.ShowAndRun()
}
