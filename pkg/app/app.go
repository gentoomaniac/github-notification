package app

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/google/go-github/v67/github"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"

	"github.com/gentoomaniac/github-notifications/pkg/gh"
)

const browserBinary = "google-chrome-stable"

var unread = map[bool]string{
	true:  "\uf52b",
	false: "",
}

var symbols = map[string]string{
	"review_requested": "\uf407",
	"author":           "\uf415",
	"state_change":     "\uf090	",
	"ci_activity":      "\ue78c",
}

func symbolLookup(s string) string {
	symbol, ok := symbols[s]
	if !ok {
		return s
	}
	return symbol
}

func New(githubToken string) App {
	return App{
		ghWrapper:    gh.New(githubToken),
		pullrequests: make(map[string]*github.PullRequest),
	}
}

type App struct {
	ghWrapper     *gh.Github
	notifications *tview.Table
	details       *tview.Form
	pullrequests  map[string]*github.PullRequest
}

func (a *App) Run() {
	app := tview.NewApplication()

	notifications, err := a.ghWrapper.GetNotifications()
	if err != nil {
		log.Error().Err(err).Msg("failed getting notifications")
	}

	a.notifications = a.Notifications(notifications)
	a.details = a.Details()
	a.updateDetails(notifications[0])

	layout := tview.NewFlex().
		AddItem(a.notifications, 0, 1, true).
		AddItem(a.details, 0, 1, false)

	// Shortcuts to navigate the slides.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.Stop()
			return nil
		} else if event.Rune() == rune('n') {
			browser(browserBinary, "https://github.com/notifications")
			return nil
		}

		return event
	})

	if err := app.SetRoot(layout, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}

func (a *App) Details() *tview.Form {
	form := tview.NewForm().
		AddTextView("Title", "", 60, 2, true, false).
		AddTextView("URL", "", 60, 2, true, false).
		AddTextView("State", "", 20, 2, true, false)
	// AddInputField("Last name", "", 20, nil, nil).
	// AddTextArea("Address", "", 40, 0, 0, nil).
	// AddTextView("Notes", "This is just a demo.\nYou can enter whatever you wish.", 40, 2, true, false).
	// AddCheckbox("Age 18+", false, nil).
	// AddPasswordField("Password", "", 10, '*', nil)
	form.SetBorder(true).SetTitle("").SetTitleAlign(tview.AlignLeft)
	return form
}

func (a *App) updateDetails(notification *github.Notification) {
	a.details.SetTitle(*notification.Repository.FullName)
	a.details.GetFormItemByLabel("Title").(*tview.TextView).SetText(strings.Join([]string{unread[*notification.Unread], *notification.Subject.Title}, " "))
	a.details.GetFormItemByLabel("URL").(*tview.TextView).SetText(getHtmlUrl(notification))
	if notification.Subject.GetType() == "PullRequest" && a.pullrequests[notification.GetID()] != nil {
		a.details.GetFormItemByLabel("URL").(*tview.TextView).SetText(*a.pullrequests[notification.GetID()].State)
	}
}

func (a *App) Notifications(notifications []*github.Notification) *tview.Table {
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
		go browser(browserBinary, getHtmlUrl(notifications[row]))
	}).SetSelectionChangedFunc(func(row int, column int) {
		a.updateDetails(notifications[row])
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == rune('r') {
			row, _ := table.GetSelection()
			n := notifications[row]

			id, err := strconv.ParseInt(n.GetID(), 10, 64)
			if err != nil {
				log.Error().Err(err).Msg("failed parsing PR id")
			}

			if n.GetSubject().GetType() == "PullRequest" {
				pr, err := a.ghWrapper.GetPr(n.GetRepository().GetOwner().GetLogin(), n.GetRepository().GetName(), int(id))
				if err != nil {
					fmt.Println(err)
				}
				a.pullrequests[n.GetID()] = pr
			}

			return nil
		}

		return event
	})

	return table
}

func browser(binary string, arg ...string) {
	cmd := exec.Command(binary, arg...)
	err := cmd.Start()
	if err != nil {
		log.Error().Err(err).Msg("failed starting browser")
	}
}

func getHtmlUrl(n *github.Notification) string {
	fields := strings.Split(n.GetSubject().GetURL(), "/")
	return fmt.Sprintf("https://github.com/%s/pull/%s", *n.Repository.FullName, fields[len(fields)-1])
}

func markRead(n *github.Notification) {}

func delete(n *github.Notification) {}
