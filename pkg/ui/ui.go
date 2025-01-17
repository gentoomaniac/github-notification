package ui

import (
	"log"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/google/go-github/v67/github"
	"github.com/rivo/tview"
)

var unread = map[bool]string{
	true:  "\uf52b",
	false: "",
}

var symbols = map[string]string{
	"review_requested": "\uf407",
	"author":           "\uf415",
}

func symbolLookup(s string) string {
	symbol, ok := symbols[s]
	if !ok {
		return s
	}
	return symbol
}

type Ui struct {
	selectedIndex int

	notifications *tview.Table
	details       *tview.Form
}

func (u *Ui) Run(notifications []*github.Notification) {
	app := tview.NewApplication()

	u.notifications = u.Notifications(notifications)
	u.details = u.Details()
	u.updateDetails(notifications[0])

	layout := tview.NewFlex().
		AddItem(u.notifications, 0, 1, true).
		AddItem(u.details, 0, 1, false)

	// Shortcuts to navigate the slides.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.Stop()
			return nil
		}
		return event
	})

	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (u *Ui) Details() *tview.Form {
	form := tview.NewForm().
		AddTextView("Title", "", 60, 2, true, false).
		AddTextView("URL", "", 60, 2, true, false)
	// AddInputField("First name", "", 20, nil, nil).
	// AddInputField("Last name", "", 20, nil, nil).
	// AddTextArea("Address", "", 40, 0, 0, nil).
	// AddTextView("Notes", "This is just a demo.\nYou can enter whatever you wish.", 40, 2, true, false).
	// AddCheckbox("Age 18+", false, nil).
	// AddPasswordField("Password", "", 10, '*', nil)
	form.SetBorder(true).SetTitle("").SetTitleAlign(tview.AlignLeft)
	return form
}

func (u *Ui) updateDetails(notification *github.Notification) {
	u.details.SetTitle(*notification.Repository.FullName)
	u.details.GetFormItemByLabel("Title").(*tview.TextView).SetText(strings.Join([]string{unread[*notification.Unread], *notification.Subject.Title}, " "))
	u.details.GetFormItemByLabel("URL").(*tview.TextView).SetText(*notification.Subject.URL)
}

func (u *Ui) Notifications(notifications []*github.Notification) *tview.Table {
	table := tview.NewTable().SetBorders(false).SetSelectable(true, false)
	table.SetTitle("Notification")
	for r := 0; r < len(notifications); r++ {
		table.SetCell(r, 1,
			tview.NewTableCell(unread[(*notifications[r].Unread)]).
				SetAlign(tview.AlignCenter))

		table.SetCell(r, 2,
			tview.NewTableCell(symbolLookup(*notifications[r].Reason)).
				SetAlign(tview.AlignCenter))

		table.SetCell(r, 3,
			tview.NewTableCell(*notifications[r].Repository.FullName).
				SetAlign(tview.AlignLeft))

		table.SetCell(r, 4,
			tview.NewTableCell(*notifications[r].Subject.Title).
				SetAlign(tview.AlignLeft))

	}
	table.Select(0, 0).SetFixed(1, 1).SetSelectedFunc(func(row int, column int) {
		// TODO: make configurable
		go browser("google-chrome-stable", *notifications[row].Subject.URL)
	}).SetSelectionChangedFunc(func(row int, column int) {
		u.updateDetails(notifications[row])
	})

	return table
}

func browser(binary string, arg ...string) {
	cmd := exec.Command(binary, arg...)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func markRead(n *github.Notification) {}

func delete(n *github.Notification) {}
