package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	gh "github.com/cli/go-gh"
	ghRepo "github.com/cli/go-gh/pkg/repository"
	"github.com/cli/go-gh/pkg/api"
	ghTableprinter "github.com/cli/go-gh/pkg/tableprinter"
	ghTerm "github.com/cli/go-gh/pkg/term"
)

type ListOptions struct {
	Repo   string
	Limit  int
}

func NewCmdList() *cobra.Command {
	f := ListOptions{}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists the issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if len(args) != 0 {
				return fmt.Errorf(fmt.Sprintf("Invalid argument(s). Expected 0 received %d", len(args)))
			}

			repo, err := GetRepo(f.Repo)
			if err != nil {
				return err
			}

			opts := api.ClientOptions{
				Host:    repo.Host(),
			}
			restClient, err := gh.RESTClient(&opts)

			pathComponent := fmt.Sprintf("repos/%s/%s/actions/runs?per_page=%d", repo.Owner(), repo.Name(), f.Limit)
			var apiResults ApiResponse
			err = restClient.Get(pathComponent, &apiResults)
			if err != nil {
				return err
			}
			PrettyPrint(apiResults.WorkflowRuns)
			return nil
		},
	}

	listCmd.Flags().StringVarP(&f.Repo, "repo", "R", "", "Select another repository for listing issues")
	listCmd.Flags().IntVarP(&f.Limit, "limit", "L", 5, "Limit of workflow runs to display")

	return listCmd
}

func PrettyPrint(workflowRuns []WorkflowRun) {
	terminal := ghTerm.FromEnv()
	w, _, _ := terminal.Size()
	tp := ghTableprinter.New(terminal.Out(), terminal.IsTerminalOutput(), w)

	for _, workflowRun := range workflowRuns {
		tp.AddField(strconv.Itoa(workflowRun.Id))
		tp.AddField(workflowRun.Name)
		tp.AddField(workflowRun.Status)
		tp.EndRow()
	}

	_ = tp.Render()
}

func GetRepo(r string) (ghRepo.Repository, error) {
	if r != "" {
		return ghRepo.Parse(r)
	}

	return gh.CurrentRepository()
}

type ApiResponse struct {
	TotalCount    int            `json:"total_count"`
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

type WorkflowRun struct {
	Id             int     `json:"id"`
	Name           string  `json:"name"`
	Status         string  `json:"status"`
}