package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"time"
)

type UIComponents struct {
	StartTime    *time.Time
	StopTime     *time.Time
	Duration     *time.Duration
	SummaryText  *widget.Label
	StartText    *widget.Label
	StopText     *widget.Label
	DurationText *widget.Label
	JiraId       *widget.Entry
}

func createTimeButtons(ui *UIComponents) *fyne.Container {
	startButton := widget.NewButton("Start", func() {
		*ui.StartTime = time.Now()
		ui.StartText.SetText(ui.StartTime.Format(time.RFC3339))
		log.Println(ui.StartTime.Format(time.RFC3339))
	})
	
	stopButton := widget.NewButton("Stop", func() {
		*ui.StopTime = time.Now()
		ui.StopText.SetText(ui.StopTime.Format(time.RFC3339))
		log.Println(ui.StopTime.Format(time.RFC3339))
		*ui.Duration = ui.StopTime.Sub(*ui.StartTime)
		ui.DurationText.SetText(ui.Duration.String())
	})
	
	resetButton := widget.NewButton("Reset", func() {
		ui.StartText.SetText("")
		ui.StopText.SetText("")
		ui.DurationText.SetText("")
	})
	
	return container.NewHBox(startButton, stopButton, resetButton)
}

func createJiraIdContainer(ui *UIComponents) *fyne.Container {
	jiraButton := widget.NewButton("pull", func() {
		log.Println("id:", ui.JiraId.Text)
		ui.SummaryText.SetText("")
		// Use jiraApiKey instead of fetching token
		jiraItem := getJiraItem(ui.JiraId.Text, jiraApiKey)
		bytJiraItem := []byte(jiraItem)
		var mapData map[string]interface{}
		json.Unmarshal(bytJiraItem, &mapData)
		log.Println(mapData)
	})
	
	jiraIdContainer := container.NewWithoutLayout(ui.JiraId, jiraButton)
	ui.JiraId.Resize(fyne.NewSize(200, 32))
	ui.JiraId.Move(fyne.NewPos(0, 0))
	ui.JiraId.Wrapping = fyne.TextWrapOff
	ui.JiraId.Scroll = container.ScrollNone
	jiraButton.Resize(fyne.NewSize(40, 32))
	jiraButton.Move(fyne.NewPos(208, 0))
	
	container := container.NewVBox()
	container.Add(widget.NewLabel("Jira Id"))
	container.Add(jiraIdContainer)
	
	return container
}

func createJiraItemContainer(ui *UIComponents) *fyne.Container {
	jiraItemContainer := container.NewVBox()
	jiraItemContainer.Add(container.NewHBox(widget.NewLabel("Summary"), ui.SummaryText))
	jiraItemContainer.Add(container.NewHBox(widget.NewLabel("Start"), ui.StartText))
	jiraItemContainer.Add(container.NewHBox(widget.NewLabel("End"), ui.StopText))
	jiraItemContainer.Add(container.NewHBox(widget.NewLabel("Duration"), ui.DurationText))
	return jiraItemContainer
}

func createMainForm(ui *UIComponents) *widget.Form {
	buttonContainer := createTimeButtons(ui)
	jiraIdContainer := createJiraIdContainer(ui)
	jiraItemContainer := createJiraItemContainer(ui)
	
	return &widget.Form{
		Items: []*widget.FormItem{
			{Text: "", Widget: jiraIdContainer},
			{Text: "", Widget: buttonContainer},
			{Text: "", Widget: jiraItemContainer},
		},
		OnSubmit: func() {
			log.Println("id:", ui.JiraId.Text)
			log.Println("start:", ui.StartTime.Format(time.RFC3339))
			log.Println("stop:", ui.StopTime.Format(time.RFC3339))
			log.Println("duration:", *ui.Duration)
		},
	}
}