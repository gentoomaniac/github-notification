package gh

import (
	"context"

	"github.com/google/go-github/v67/github"
)

const PullRequest = "PullRequest"

const (
	ReviewRequested = "review_requested"
	Author          = "author"
	StateChange     = "state_change"
	CiActivity      = "ci_activity"
)

func New() *Github {
	return &Github{
		context: context.Background(),
	}
}

type Github struct {
	context       context.Context
	notifications []*github.Notification
}

func (g *Github) GetNotifications(notificationToken string) ([]*github.Notification, error) {
	client := github.NewClient(nil).WithAuthToken(notificationToken)
	options := &github.NotificationListOptions{
		All: false,
		ListOptions: github.ListOptions{
			Page: 0,
		},
	}

	n, res, err := client.Activity.ListNotifications(
		g.context,
		options,
	)
	if err != nil {
		return g.notifications, err
	}
	g.notifications = n

	for i := 1; i < res.LastPage; i++ {
		n, res, err = client.Activity.ListNotifications(
			g.context,
			options,
		)
		if err != nil {
			return g.notifications, err
		}
		g.notifications = append(g.notifications, n...)
	}

	return g.notifications, nil
}

func (g *Github) GetPr(orgToken string, owner string, repo string, id int) (*github.PullRequest, error) {
	client := github.NewClient(nil).WithAuthToken(orgToken)
	pr, _, err := client.PullRequests.Get(context.Background(), owner, repo, id)

	return pr, err
}
