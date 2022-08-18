package cmd

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/cli/go-gh/pkg/api"
	"gopkg.in/h2non/gock.v1"
)

func TestListWithIncorrectArguments(t *testing.T) {
	t.Cleanup(gock.Off)
	cmd := NewCmdList()
	cmd.SetArgs([]string{"INCORRECT"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("Invalid argument(s). Expected 0 received 1"))
}

func TestListWithSuccess(t *testing.T) {
	t.Cleanup(gock.Off)
	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/runs").
		Reply(200).
		JSON(`{
			"total_count":1,
			"workflow_runs":[
				{
					"id":1,
					"name":"Test Workflow",
					"status":"queued"
				}
			]
		}`)

	cmd := NewCmdList()
	cmd.SetArgs([]string{"-R", "testOrg/testRepo"})
	err := cmd.Execute()

	assert.Nil(t, err)

	assert.True(t, gock.IsDone(), PrintPendingMocks(gock.Pending()))
}

func TestListWithInternalServerError(t *testing.T) {
	t.Cleanup(gock.Off)
	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/runs").
		Reply(500).
		JSON(`{
			"message": "Internal Server Error",
			"documentation_url": "https://docs.github.com/rest/reference/actions#get-github-actions-cache-list-for-a-repository"
		}`)

	cmd := NewCmdList()
	cmd.SetArgs([]string{"-R", "testOrg/testRepo"})
	err := cmd.Execute()

	assert.NotNil(t, err)

	var httpError api.HTTPError
	errors.As(err, &httpError)
	assert.Equal(t, httpError.StatusCode, 500)

	assert.True(t, gock.IsDone(), PrintPendingMocks(gock.Pending()))
}

func TestListWithPendingMocks(t *testing.T) {
	t.Cleanup(gock.Off)
	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/runs/xyz").
		Reply(500).
		JSON(`{
			"message": "Internal Server Error",
			"documentation_url": "https://docs.github.com/rest/reference/actions#get-github-actions-cache-list-for-a-repository"
		}`)

	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/runs").
		Reply(500).
		JSON(`{
			"message": "Internal Server Error",
			"documentation_url": "https://docs.github.com/rest/reference/actions#get-github-actions-cache-list-for-a-repository"
		}`)

	cmd := NewCmdList()
	cmd.SetArgs([]string{"-R", "testOrg/testRepo"})
	err := cmd.Execute()

	assert.NotNil(t, err)

	assert.False(t, gock.IsDone(), PrintPendingMocks(gock.Pending()))
}

func PrintPendingMocks(mocks []gock.Mock) string {
	paths := []string{}
	for _, mock := range mocks {
		paths = append(paths, mock.Request().URLStruct.String())
	}
	return fmt.Sprintf("%d unmatched mocks: %s", len(paths), strings.Join(paths, ", "))
}
