package gh

import (
	"context"

	"github.com/google/go-github/v67/github"
)

func New(token string) *Github {
	return &Github{
		client:  github.NewClient(nil).WithAuthToken(token),
		context: context.Background(),
	}
}

type Github struct {
	client        *github.Client
	context       context.Context
	notifications []*github.Notification
}

func (g *Github) GetNotifications() ([]*github.Notification, error) {
	options := &github.NotificationListOptions{
		All: false,
		ListOptions: github.ListOptions{
			Page: 0,
		},
	}

	n, res, err := g.client.Activity.ListNotifications(
		g.context,
		options,
	)
	if err != nil {
		return g.notifications, err
	}
	g.notifications = n

	for i := 1; i < res.LastPage; i++ {
		n, res, err = g.client.Activity.ListNotifications(
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

func (g *Github) GetPr(owner string, repo string, id int) (*github.PullRequest, error) {
	pr, _, err := g.client.PullRequests.Get(g.context, owner, repo, id)

	return pr, err
}
