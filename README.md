
# gh-cli-extension-demo

This repo is part of a demo where we learn how to create an extension for gh-cli in golang.

## Libraries Used
1. [go-gh](https://github.com/cli/go-gh) - A Go module for interacting with gh and the GitHub API from the command line.
2. [Cobra](https://github.com/spf13/cobra) - It is a library used giving structure and managing inputs in our cli


## Development Setup

1. Install [golang](https://go.dev/doc/install)
2. Install the `gh` CLI - see the [installation](https://github.com/cli/cli#installation)
   
   _Installation requires a minimum version (2.0.0) of the the GitHub CLI that supports extensions._

## Stage 1

In this stage we will:-

    1. Create a basic extension 
    2. properly structure the folder structure
    3. run our extension

**A. Create extension** 
		
		gh extension create --precompiled=go <extension-name-excluding-gh>
	
This command will generate a new precompiled golang extension. 

**B. Install cobra**

		go get -u github.com/spf13/cobra@latest
	
**C. Update directory structure**
While we are welcome to provide our own organization, typically a Cobra-based application will follow the following organizational structure:

	  ▾ appName/
	    ▾ cmd/
	        root.go
	      main.go
In a Cobra app, typically the main.go file is very bare. It serves one purpose: initializing Cobra.

	package main

	import (
	  "{pathToYourApp}/cmd"
	)

	func main() {
	  cmd.Execute()
	}
	
**D. Final Code** 

root.go

	package cmd

	import (
		"fmt"
		"os"

		"github.com/spf13/cobra"
	)

	var  rootCmd  =  &cobra.Command{
		Use: "gh cli-extension-demo",
		Short: "This extension prints out all the input user provides",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Stage 1 of demo is done")
		},
	}

	func  Execute() {
		if  err  := rootCmd.Execute(); err !=  nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

main.go

	package main

	import  "github.com/t-dedah/gh-cli-extension-demo/cmd"

	func  main() {
		cmd.Execute()
	}


**E. Build and Test**

1. Build your project using `go build`. This will generate a new binary for the cli. 

2. Install the binary using `gh extension install .`. This will create a symlink between the binary this repo and go package directory in your machine. Any further builds will not require the installation step.

3. Test the extension using `gh <extension name>` This should print `Stage 1 of demo is done` on the terminal

## Stage 2

In this stage we will:-

1. Add a new command
2. Take input for command
3. Make api calls

**A. Add new command**
1. Create a new file `list.go` on cmd folder
2. Add this to `list.go`

```
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	gh "github.com/cli/go-gh"
	ghRepo "github.com/cli/go-gh/pkg/repository"
)

type ListOptions struct {
	Repo   string
}

func NewCmdList() *cobra.Command {
	f := ListOptions{}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists the issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Input Repo: %s", f.Repo)
			return nil
		},
	}

	listCmd.Flags().StringVarP(&f.Repo, "repo", "R", "", "Select another repository for listing issues")

	return listCmd
}
```

This will take one string argument `-R` and print the value for that input.

**B. Add list command to root**
1. Add new function to `root.go`

```
func addCommandsToRoot() {
	rootCmd.AddCommand(NewCmdList())
}
```

2. Call this function from `Execute()`

```
func Execute() {
	addCommandsToRoot()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

3. Build and test the extension

**C. Use go-gh to make api calls**
1. Update the list.go 

```
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	gh "github.com/cli/go-gh"
	ghRepo "github.com/cli/go-gh/pkg/repository"
	"github.com/cli/go-gh/pkg/api"
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
			fmt.Println(apiResults.WorkflowRuns)
			return nil
		},
	}

	listCmd.Flags().StringVarP(&f.Repo, "repo", "R", "", "Select another repository for listing issues")
	listCmd.Flags().IntVarP(&f.Limit, "limit", "L", 5, "Limit of workflow runs to display")

	return listCmd
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
```

**D. Add table printer**

1. This function will tabulated data using go-gh table printer

```
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
```

2. Final list.go

```
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
```

3. Build and test
