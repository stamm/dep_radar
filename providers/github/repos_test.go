package github

import (
	"context"
	"errors"
	"testing"

	"github.com/stamm/dep_radar/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Repos

func TestGithubReposOrg_Ok(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &mocks.IWebClient{}
	client.On("Get", mock.Anything, "https://api.github.com/orgs/dep-radar/repos").Return([]byte(`[
  {
    "full_name": "dep-radar/test_app"
  },
  {
    "full_name": "dep-radar/test_app2"
  }
]`), nil)
	tagsGetter := New(client)

	pkgs, err := tagsGetter.GetAllOrgRepos(context.Background(), "dep-radar")
	require.NoError(err)
	require.Len(pkgs, 2, "Expect 2 repos in org")
	require.EqualValues("github.com/dep-radar/test_app", pkgs[0])
	require.EqualValues("github.com/dep-radar/test_app2", pkgs[1])
}

func TestGithubReposOrg_Error(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &mocks.IWebClient{}
	client.On("Get", mock.Anything, "https://api.github.com/orgs/dep-radar/repos").Return([]byte(``), errors.New("aaa"))
	tagsGetter := New(client)

	pkgs, err := tagsGetter.GetAllOrgRepos(context.Background(), "dep-radar")
	require.Error(err)
	require.Len(pkgs, 0, "Expect 0 repos in org")
}

func TestGithubUserExists_Ok(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &mocks.IWebClient{}
	client.On("Get", mock.Anything, "https://api.github.com/users/stamm").Return([]byte(`{}`), nil)
	tagsGetter := New(client)

	exists := tagsGetter.UserExists(context.Background(), "stamm")
	require.True(exists, "Expect user exists")
}

func TestGithubUserExists_Error(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &mocks.IWebClient{}
	client.On("Get", mock.Anything, "https://api.github.com/users/stamm").Return(nil, errors.New("err"))
	tagsGetter := New(client)

	exists := tagsGetter.UserExists(context.Background(), "stamm")
	require.False(exists, "Expect user doesn't exist")
}

func TestGithubReposUser_Ok(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &mocks.IWebClient{}
	client.On("Get", mock.Anything, "https://api.github.com/users/stamm/repos").Return([]byte(`[
  {
    "full_name": "stamm/dep_radar"
  },
  {
    "full_name": "stamm/callstat"
  }
]`), nil)
	tagsGetter := New(client)

	pkgs, err := tagsGetter.GetAllUserRepos(context.Background(), "stamm")
	require.NoError(err)
	require.Len(pkgs, 2, "Expect 2 repos for user")
	require.EqualValues("github.com/stamm/dep_radar", pkgs[0])
	require.EqualValues("github.com/stamm/callstat", pkgs[1])
}

func TestGithubReposUser_Error(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &mocks.IWebClient{}
	client.On("Get", mock.Anything, "https://api.github.com/users/stamm/repos").Return([]byte(``), errors.New("aaa"))
	tagsGetter := New(client)

	pkgs, err := tagsGetter.GetAllUserRepos(context.Background(), "stamm")
	require.Error(err)
	require.Len(pkgs, 0, "Expect 0 repos for user")
}
