package ghch

import "github.com/octokit/go-octokit/octokit"

func reducePR(pr *octokit.PullRequest) *octokit.PullRequest {
	return &octokit.PullRequest{
		HTMLURL:        pr.HTMLURL,
		Title:          pr.Title,
		Number:         pr.Number,
		State:          pr.State,
		Body:           pr.Body,
		CreatedAt:      pr.CreatedAt,
		UpdatedAt:      pr.UpdatedAt,
		MergedAt:       pr.MergedAt,
		MergeCommitSha: pr.MergeCommitSha,
		User:           reduceUser(pr.User),
		Head:           reducePullRequestCommit(pr.Head),
		Base:           reducePullRequestCommit(pr.Base),
		MergedBy:       reduceUser(pr.MergedBy),
	}
}

func reduceUser(u octokit.User) octokit.User {
	return octokit.User{
		Login:     u.Login,
		AvatarURL: u.AvatarURL,
		Type:      u.Type,
		HTMLURL:   u.HTMLURL,
	}
}

func reduceRepo(r *octokit.Repository) *octokit.Repository {
	if r == nil {
		return nil
	}
	return &octokit.Repository{
		Owner:    reduceUser(r.Owner),
		Name:     r.Name,
		FullName: r.FullName,
		HTMLURL:  r.HTMLURL,
	}
}

func reducePullRequestCommit(prc octokit.PullRequestCommit) octokit.PullRequestCommit {
	return octokit.PullRequestCommit{
		Label: prc.Label,
		Ref:   prc.Ref,
		Sha:   prc.Sha,
		User:  reduceUser(prc.User),
		Repo:  reduceRepo(prc.Repo),
	}
}
