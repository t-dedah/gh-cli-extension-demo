
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
