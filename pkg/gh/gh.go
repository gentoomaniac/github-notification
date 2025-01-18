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

func New(classicToken string, finegrainedToken string) *Github {
	return &Github{
		classicClient:     github.NewClient(nil).WithAuthToken(classicToken),
		finegrainedClient: github.NewClient(nil).WithAuthToken(finegrainedToken),
		context:           context.Background(),
	}
}

type Github struct {
	classicClient     *github.Client
	finegrainedClient *github.Client
	context           context.Context
	notifications     []*github.Notification
}

func (g *Github) GetNotifications() ([]*github.Notification, error) {
	options := &github.NotificationListOptions{
		All: false,
		ListOptions: github.ListOptions{
			Page: 0,
		},
	}

	n, res, err := g.classicClient.Activity.ListNotifications(
		g.context,
		options,
	)
	if err != nil {
		return g.notifications, err
	}
	g.notifications = n

	for i := 1; i < res.LastPage; i++ {
		n, res, err = g.classicClient.Activity.ListNotifications(
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
	pr, _, err := g.finegrainedClient.PullRequests.Get(context.Background(), owner, repo, id)

	return pr, err
}
