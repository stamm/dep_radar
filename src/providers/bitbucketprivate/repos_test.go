package bitbucketprivate

import (
	"context"
	"errors"
	"testing"

	"github.com/stamm/dep_radar/interfaces/mocks"
	"github.com/stretchr/testify/require"
)

// Repos

func TestBBRepos_Ok(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://bitbucket.example.com/rest/api/1.0/projects/proj/repos?start=0").Return([]byte(`{
	"isLastPage": true,
	"values": [
	{
		"slug": "test"
	},
	{
		"slug": "test2"
	}
	]
}`), nil)
	provider := New(mHttpClient, "bitbucket.example.com", "godep.example.com", "https://bitbucket.example.com")

	pkgs, err := provider.GetAllRepos(context.Background(), "proj")
	require.NoError(err)
	require.Len(pkgs, 2, "Expect 2 repos")
	require.EqualValues("godep.example.com/test", pkgs[0])
	require.EqualValues("godep.example.com/test2", pkgs[1])
}

func TestBBRepos_TwoPages(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://bitbucket.example.com/rest/api/1.0/projects/proj/repos?start=0").Return([]byte(`{
	"isLastPage": false,
	"values": [
	{ "slug": "test" }
	],
	"nextPageStart": 1
}`), nil)
	mHttpClient.On("Get", "https://bitbucket.example.com/rest/api/1.0/projects/proj/repos?start=1").Return([]byte(`{
	"isLastPage": true,
	"values": [
	{ "slug": "test2" }
	]
}`), nil)
	provider := New(mHttpClient, "bitbucket.example.com", "godep.example.com", "https://bitbucket.example.com")

	pkgs, err := provider.GetAllRepos(context.Background(), "proj")
	require.NoError(err)
	require.Len(pkgs, 2, "Expect 2 repos")
	require.EqualValues("godep.example.com/test", pkgs[0])
	require.EqualValues("godep.example.com/test2", pkgs[1])
}

func TestGithubRepos_Error(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://bitbucket.example.com/rest/api/1.0/projects/proj/repos?start=0").Return(nil, errors.New("error"))

	provider := New(mHttpClient, "bitbucket.example.com", "godep.example.com", "https://bitbucket.example.com")

	pkgs, err := provider.GetAllRepos(context.Background(), "proj")
	require.Error(err)
	require.Len(pkgs, 0, "Expect 0 repos")
}
