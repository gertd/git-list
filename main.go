package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/gertd/git-list/version"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

// command line args
const (
	GitHost  = "host"
	GitUser  = "user"
	GitOrg   = "org"
	GitToken = "token"
	Verbose  = "verbose"
)

// environment variables
const (
	GitHostEnv  = "GIT_HOST"
	GitUserEnv  = "GIT_USER"
	GitOrgEnv   = "GIT_ORG"
	GitTokenEnv = "GIT_TOKEN"
)

// command line usage
const (
	GitHostUsage  = "GitHub Enterprise host address"
	GitUserUsage  = "GitHub user"
	GitOrgUsage   = "GitHub organization"
	GitTokenUsage = "GitHub private access token"
	VerboseUsage  = "verbose output"
)

// command literals
const (
	AppName  = "git-list"
	AppUsage = "list all GitHub repos"
)

func main() {
	app := cli.NewApp()
	app.Name = AppName
	app.Usage = AppUsage
	app.Action = gitList
	app.Version = version.Info()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     GitUser,
			Value:    "",
			Usage:    GitUserUsage,
			EnvVar:   GitUserEnv,
			Required: true,
		},
		cli.StringFlag{
			Name:   GitOrg,
			Value:  "",
			Usage:  GitOrgUsage,
			EnvVar: GitOrgEnv,
			// Required: true,
		},
		cli.StringFlag{
			Name:  "type",
			Value: "all",
			Usage: "all | owner | public | private | member",
		},
		cli.StringFlag{
			Name:     GitToken,
			Value:    "",
			Usage:    GitTokenUsage,
			EnvVar:   GitTokenEnv,
			Required: true,
		},
		cli.StringFlag{
			Name:   GitHost,
			Value:  "",
			Usage:  GitHostUsage,
			EnvVar: GitHostEnv,
		},
		cli.BoolFlag{
			Name:  Verbose,
			Usage: VerboseUsage,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}

	os.Exit(1)
}

func gitList(c *cli.Context) error { //nolint:funlen
	var (
		err    error
		ctx    = context.Background()
		client *github.Client
		repos  []*github.Repository
		resp   *github.Response
	)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	gitToken := c.GlobalString(GitToken)
	gitOrg := c.GlobalString(GitOrg)
	gitRepoType := c.GlobalString("type")
	gitUser := c.GlobalString(GitUser)
	gitHost := c.GlobalString(GitHost)
	verbose := c.GlobalBool(Verbose)

	if verbose {
		fmt.Printf("user: [%s] org: [%s]\n", gitUser, gitOrg)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	if gitHost == "" {
		client = github.NewClient(tc)
	} else if client, err = github.NewEnterpriseClient(gitHost, gitHost, tc); err != nil {
		return err
	}

	page := 1

	fmt.Fprintf(w, "%s\t%s\t%s\n",
		"last update",
		"repo name",
		"description")

	for {
		if len(gitOrg) == 0 {
			opt := &github.RepositoryListOptions{
				Visibility: gitRepoType,
				// Type:        gitRepoType,
				ListOptions: github.ListOptions{PerPage: 10, Page: page},
			}

			repos, resp, err = client.Repositories.List(ctx, gitUser, opt)
			if err != nil {
				return err
			}
		} else {
			opt := &github.RepositoryListByOrgOptions{
				Type:        gitRepoType,
				ListOptions: github.ListOptions{PerPage: 10, Page: page},
			}

			repos, resp, err = client.Repositories.ListByOrg(ctx, c.GlobalString(GitOrg), opt)
			if err != nil {
				return err
			}
		}

		for _, repo := range repos {
			fmt.Fprintf(w, "%04d-%02d-%02d\t%s\t%s\n",
				repo.UpdatedAt.Year(),
				repo.UpdatedAt.Month(),
				repo.UpdatedAt.Day(),
				ifNilEmpty(repo.FullName),
				ifNilEmpty(repo.Description))
		}

		if resp.NextPage == 0 {
			break
		}

		page = resp.NextPage
	}
	w.Flush()

	return nil
}

func ifNilEmpty(pstr *string) string {
	if pstr == nil {
		return ""
	}

	return *pstr
}
