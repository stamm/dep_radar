package bitbucketprivate

import (
	"context"
	"errors"
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/interfaces/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGoImport_GetProject(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", mock.Anything, "https://godep.example.com/app?go-get=1").Return([]byte(`<html>
<head>
<meta name="go-import" content="godep.example.com/app git ssh://git@bitbucket.example.com/go_project/app.git"/></head>
<body>

</body>
</html>`), nil)

	project, err := GetProject(context.Background(), mHttpClient, i.Pkg("godep.example.com/app"), "bitbucket.example.com")
	require.NoError(err)
	require.Equal("go_project", project)
}

func TestGoImport_GetProjectWithError(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", mock.Anything, "https://godep.example.com/app?go-get=1").Return(nil, errors.New("err"))

	project, err := GetProject(context.Background(), mHttpClient, i.Pkg("godep.example.com/app"), "bitbucket.example.com")
	require.EqualError(err, "err")
	require.Equal("", project)
}

func TestGoImport_GetProjectWithBadPrefix(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", mock.Anything, "https://godep.example.com/app?go-get=1").Return([]byte(`<html>
<head>
<meta name="go-import" content="godep.example.com/app git ssh://git@bitbu_cket.example.com/go_project/app.git"/></head>
<body>

</body>
</html>`), nil)

	project, err := GetProject(context.Background(), mHttpClient, i.Pkg("godep.example.com/app"), "bitbucket.example.com")
	require.EqualError(err, "Can't find project in repo urls ssh://git@bitbu_cket.example.com/go_project/app.git")
	require.Equal("", project)
}

func TestGoImport_GetProjectFromSecond(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", mock.Anything, "https://godep.example.com/app?go-get=1").Return([]byte(`<html>
<head>
<meta name="go-import" content="godep.example.com/app git ssh://git@bit_bucket.example.com/go_project/app.git"/>
<meta name="go-import" content="godep.example.com/app git ssh://git@bitbucket.example.com/go_project/app.git"/>
</head>
<body>

</body>
</html>`), nil)

	project, err := GetProject(context.Background(), mHttpClient, i.Pkg("godep.example.com/app"), "bitbucket.example.com")
	require.NoError(err)
	require.Equal("go_project", project)
}

func TestGoImport_GetProjectWrongPackage(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", mock.Anything, "https://godep.example.com/app?go-get=1").Return([]byte(`<html>
<head>
<meta name="go-import" content="godep.example.com/app2 git ssh://git@bitbucket.example.com/go_project/app.git"/>
</head>
<body>

</body>
</html>`), nil)

	project, err := GetProject(context.Background(), mHttpClient, i.Pkg("godep.example.com/app"), "bitbucket.example.com")
	require.EqualError(err, "Can't find project in repo urls ")
	require.Equal("", project)
}
