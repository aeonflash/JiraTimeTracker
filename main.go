package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"time"
)

var jiraGraphQlBaseUri = "https://dish.atlassian.net/gateway/api/graphql"
var jiraCloudId ="REDACTED_CLOUD_ID"

var jiraApiKey string



func main() {
	loadJiraConfig()
	a := app.New()
	w := a.NewWindow("JiraWidgetLite")
	w.Resize(fyne.NewSize(400, 400))
	w.SetTitle("JiraWidgetLite")

	startTime := time.Now()
	stopTime := time.Now()
	duration := time.Duration(0)

	ui := &UIComponents{
		StartTime:    &startTime,
		StopTime:     &stopTime,
		Duration:     &duration,
		SummaryText:  widget.NewLabel(""),
		StartText:    widget.NewLabel(""),
		StopText:     widget.NewLabel(""),
		DurationText: widget.NewLabel(""),
		JiraId:       widget.NewEntry(),
	}

	form := createMainForm(ui)
	w.SetContent(form)
	w.ShowAndRun()
}
