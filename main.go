package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// command line args
const (
	GitUser  = "user"
	GitOrg   = "org"
	GitToken = "token"
	Verbose  = "verbose"
)

// environment variables
const (
	GitUserEnv  = "GIT_USER"
	GitOrgEnv   = "GIT_ORG"
	GitTokenEnv = "GIT_TOKEN"
)

// command line usage
const (
	GitUserUsage  = "GitHub user"
	GitOrgUsage   = "GitHub organization"
	GitTokenUsage = "GitHub private access token"
	VerboseUsage  = "verbose output"
)

// command literals
const (
	AppName  = "git-list"
	AppUsage = "list all GitHub repos"
	GitCmd   = "git"
	GitClone = "clone"
)

func main() {

	app := cli.NewApp()
	app.Name = AppName
	app.Usage = AppUsage
	app.Action = gitList
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   GitUser,
			Value:  "",
			Usage:  GitUserUsage,
			EnvVar: GitUserEnv,
		},
		cli.StringFlag{
			Name:   GitOrg,
			Value:  "",
			Usage:  GitOrgUsage,
			EnvVar: GitOrgEnv,
		},
		cli.StringFlag{
			Name:   GitToken,
			Value:  "",
			Usage:  GitTokenUsage,
			EnvVar: GitTokenEnv,
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

func gitList(c *cli.Context) error {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	gitToken := c.GlobalString(GitToken)
	gitOrg := c.GlobalString(GitOrg)
	gitUser := c.GlobalString(GitUser)
	verbose := c.GlobalBool(Verbose)

	if verbose {
		fmt.Printf("user: [%s] org: [%s]\n", gitUser, gitOrg)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	page := 1

	var repos []*github.Repository
	var resp *github.Response
	var err error

	fmt.Fprintf(w, "%s\t%s\t%s\n",
		"last update",
		"repo name",
		"description")

	for {

		if len(gitOrg) == 0 {

			opt := &github.RepositoryListOptions{
				ListOptions: github.ListOptions{PerPage: 10, Page: page},
			}

			repos, resp, err = client.Repositories.List(gitUser, opt)
			if err != nil {
				return err
			}

		} else {
			opt := &github.RepositoryListByOrgOptions{
				ListOptions: github.ListOptions{PerPage: 10, Page: page},
			}

			repos, resp, err = client.Repositories.ListByOrg(c.GlobalString(GitOrg), opt)
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
