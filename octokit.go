package ghch

import (
	"github.com/google/go-github/github"
)

func reducePR(pr *github.PullRequest) *github.PullRequest {
	return &github.PullRequest{
		HTMLURL:        pr.HTMLURL,
		Title:          pr.Title,
		Number:         pr.Number,
		State:          pr.State,
		Body:           pr.Body,
		CreatedAt:      pr.CreatedAt,
		UpdatedAt:      pr.UpdatedAt,
		MergedAt:       pr.MergedAt,
		MergeCommitSHA: pr.MergeCommitSHA,
		User:           reduceUser(pr.User),
		Head:           reducePullRequestBranch(pr.Head),
		Base:           reducePullRequestBranch(pr.Base),
		MergedBy:       reduceUser(pr.MergedBy),
	}
}

func reduceUser(u *github.User) *github.User {
	return &github.User{
		Login:     u.Login,
		AvatarURL: u.AvatarURL,
		Type:      u.Type,
		HTMLURL:   u.HTMLURL,
	}
}

func reduceRepo(r *github.Repository) *github.Repository {
	if r == nil {
		return nil
	}
	return &github.Repository{
		Owner:    reduceUser(r.Owner),
		Name:     r.Name,
		FullName: r.FullName,
		HTMLURL:  r.HTMLURL,
	}
}

func reducePullRequestBranch(prb *github.PullRequestBranch) *github.PullRequestBranch {
	return &github.PullRequestBranch{
		Label: prb.Label,
		Ref:   prb.Ref,
		SHA:   prb.SHA,
		User:  reduceUser(prb.User),
		Repo:  reduceRepo(prb.Repo),
	}
}
