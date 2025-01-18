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
	gh.ReviewRequested: "\uf407",
	gh.Author:          "\uf415",
	gh.StateChange:     "\uf090	",
	gh.CiActivity:      "\ue78c",
}

func symbolLookup(s string) string {
	symbol, ok := symbols[s]
	if !ok {
		return s
	}
	return symbol
}

func New(classicToken string, finrgrainedToken string) App {
	return App{
		ghWrapper:    gh.New(classicToken, finrgrainedToken),
		pullrequests: make(map[string]*github.PullRequest),
	}
}

type App struct {
	ghWrapper     *gh.Github
	layout        *tview.Flex
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

	a.layout = tview.NewFlex().
		AddItem(a.notifications, 0, 1, true)

	// Shortcuts to navigate the slides.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == rune('q') {
			app.Stop()
			return nil
		} else if event.Rune() == rune('n') {
			browser(browserBinary, "https://github.com/notifications")
			return nil
		} else if event.Key() == tcell.KeyEsc {
			a.layout.RemoveItem(a.details)
		}

		return event
	})

	if err := app.SetRoot(a.layout, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}

func (a *App) PullrequestDetails(n *github.Notification, pr *github.PullRequest) *tview.Form {
	form := tview.NewForm().
		AddTextView("State", pr.GetState(), 20, 2, true, false).
		AddTextView("Body", pr.GetBody(), 60, 5, true, true).
		AddTextView("URL", getHtmlUrl(n), 60, 2, true, false)

	// AddInputField("Last name", "", 20, nil, nil).
	// AddTextArea("Address", "", 40, 0, 0, nil).
	// AddTextView("Notes", "This is just a demo.\nYou can enter whatever you wish.", 40, 2, true, false).
	// AddCheckbox("Age 18+", false, nil).
	// AddPasswordField("Password", "", 10, '*', nil)
	form.SetBorder(true).SetTitle(fmt.Sprintf("#%d %s", pr.GetID(), pr.GetTitle())).SetTitleAlign(tview.AlignLeft)
	return form
}

func (a *App) Notifications(notifications []*github.Notification) *tview.Table {
	table := tview.NewTable().SetBorders(false).SetSelectable(true, false)
	table.SetTitle("Notification")
	table.SetBorderPadding(2, 2, 2, 2)
	for r := 0; r < len(notifications); r++ {
		table.SetCell(r, 1,
			tview.NewTableCell(unread[(*notifications[r].Unread)]).
				SetAlign(tview.AlignCenter))

		table.SetCell(r, 2,
			tview.NewTableCell(symbolLookup(*notifications[r].Reason)).
				SetAlign(tview.AlignCenter))

		table.SetCell(r, 3,
			tview.NewTableCell(*notifications[r].Repository.FullName).
				SetAlign(tview.AlignLeft).SetMaxWidth(25))

		table.SetCell(r, 4,
			tview.NewTableCell(*notifications[r].Subject.Title).
				SetAlign(tview.AlignLeft))

	}
	table.Select(0, 0).SetFixed(1, 1).SetSelectedFunc(func(row int, column int) {
		n := notifications[row]

		if n.Subject.GetType() == gh.PullRequest {
			subject_id, err := GetIdFromNotification(n)
			if err != nil {
				log.Error().Err(err).Msg("")
			}
			id := fmt.Sprintf("%s%d", n.GetRepository().GetFullName(), subject_id)
			_, ok := a.pullrequests[id]
			if !ok {
				pr, _ := a.ghWrapper.GetPr(n.GetRepository().Owner.GetLogin(), n.GetRepository().GetName(), subject_id)
				a.pullrequests[id] = pr
			}
			if a.pullrequests[id] != nil {

				if a.layout.GetItemCount() > 1 {
					a.layout.RemoveItem(a.details)
				}
				a.details = a.PullrequestDetails(n, a.pullrequests[id])
				a.layout.AddItem(a.details, 0, 1, true)
			}
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == rune('o') {
			row, _ := table.GetSelection()
			// TODO: make configurable
			go browser(browserBinary, getHtmlUrl(notifications[row]))
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

func GetIdFromNotification(n *github.Notification) (int, error) {
	if n.GetSubject().GetType() == gh.PullRequest {
		fields := strings.Split(n.GetSubject().GetURL(), "/")
		id, err := strconv.ParseInt(fields[len(fields)-1], 10, 64)
		return int(id), err
	}
	return -1, fmt.Errorf("unknown notification type")
}

func getHtmlUrl(n *github.Notification) string {
	fields := strings.Split(n.GetSubject().GetURL(), "/")
	return fmt.Sprintf("https://github.com/%s/pull/%s", *n.Repository.FullName, fields[len(fields)-1])
}

func markRead(n *github.Notification) {}

func delete(n *github.Notification) {}
