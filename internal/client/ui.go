package client

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Elements of tview packages.
type UI struct {
	app *tview.Application
	// Need to keep this TextView to update the list of messages.
	messages *tview.TextView
}

// Run the application.
func (ui *UI) Start() {
	err := ui.app.Run()
	if err != nil {
		log.Fatalf("| Error while trying to start the client UI %v", err)
	}
}

// Append a new message in the list.
//
// {msg} is the message to append at the end of the list.
func (ui *UI) Append(msg Message) {
	name := msg.Name
	if msg.Type != NewMessage {
		name = ""
	}

	pad := 15 - len(name)
	if pad < 0 {
		pad = 0
	}

	now := time.Now()
	hour := fmt.Sprint(now.Hour()) + ":"
	mins := fmt.Sprint(now.Minute())
	if len(mins) == 1 {
		mins = "0" + mins
	}

	prefix := strings.Repeat(" ", pad) + name + " | "

	content := ""
	switch msg.Type {
	case NewMessage:
		content = msg.Content
	case Connect:
		content = msg.Name + " joined the room."
	case Disconnect:
		content = msg.Name + " leaved the room."
	}

	line := hour + mins + prefix + content
	if msg.Type != NewMessage {
		line = "[grey::i]" + line + "[-::-]"
	}

	_, err := fmt.Fprintf(ui.messages, "%s\n", line)
	if err != nil {
		log.Fatalf("| Error while trying to append msg in UI %v", err)
	}
	ui.app.Draw()
}

// Create all stuff for the UI (app, chat list, text input).
// And returns the created UI struct.
//
// {inputDoneFunc} callback when the text input is validated.
func NewUI(name string, inputDoneFunc func(value string)) *UI {
	messages := newChatView(name)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(messages, 0, 1, false).
		AddItem(newInputField(inputDoneFunc), 3, 0, false)

	app := tview.NewApplication().
		SetRoot(flex, true).
		EnableMouse(true)

	return &UI{
		app:      app,
		messages: messages,
	}
}

func newChatView(name string) *tview.TextView {
	chat := tview.NewTextView().SetDynamicColors(true)

	chat.SetBorder(true).
		SetTitle("| Chat ("+name+") |").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 2, 2)

	return chat
}

func newInputField(doneFunc func(value string)) *tview.InputField {
	inputField := tview.NewInputField().
		SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			doneFunc(inputField.GetText())
			inputField.SetText("")
		}
	})

	inputField.SetBorder(true).
		SetBackgroundColor(tcell.ColorBlack).
		SetTitleAlign(tview.AlignLeft).
		SetTitle("| Enter your message |").
		SetBorderPadding(0, 0, 2, 2)

	return inputField
}
