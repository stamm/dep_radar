package goimport

import (
	"context"
	"errors"
	"testing"

	"github.com/stamm/dep_radar/src/interfaces/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Repos

func TestGetImports_Ok(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &mocks.IWebClient{}
	client.On("Get", mock.Anything, "https://test?go-get=1").Return([]byte(`<meta name="go-import" content="prefix vcs reporoot" />`), nil)
	meta, err := GetImports(context.Background(), client, "test")

	require.NoError(err)
	require.Len(meta, 1, "Expect 1")
	require.EqualValues("prefix", meta[0].Prefix)
	require.EqualValues("vcs", meta[0].VCS)
	require.EqualValues("reporoot", meta[0].RepoRoot)
}

func TestGetImports_Error(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &mocks.IWebClient{}
	client.On("Get", mock.Anything, "https://test?go-get=1").Return([]byte{}, errors.New("err"))
	meta, err := GetImports(context.Background(), client, "test")

	require.Error(err)
	require.Len(meta, 0)
}

func TestGetImports_ErrorBody(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &mocks.IWebClient{}
	client.On("Get", mock.Anything, "https://test?go-get=1").Return([]byte(`<meta name="go-import" content="prefix vcs reporoot"`), nil)
	meta, err := GetImports(context.Background(), client, "test")

	require.Error(err)
	require.Len(meta, 0)
}
