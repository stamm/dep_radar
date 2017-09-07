package bitbucketprivate

import (
	"errors"
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/interfaces/mocks"
	"github.com/stretchr/testify/require"
)

func TestBBRepo_SetProject(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://godep.example.com/app?go-get=1").Return([]byte(`<html>
<head>
<meta name="go-import" content="godep.example.com/app git ssh://git@bitbucket.example.com/go_project/app.git"/></head>
<body>

</body>
</html>`), nil)

	prov := New(mHttpClient, "bitbucket.example.com", "godep.example.com", "https://bitbucket.example.com")
	project, err := prov.getProject(i.Pkg("godep.example.com/app"))
	require.NoError(err)
	require.Equal("go_project", project)
}

func TestBBRepo_SetProject_DontSetTwice(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	prov := New(&mocks.IWebClient{}, "bitbucket.example.com", "godep.example.com", "https://bitbucket.example.com")
	prov.mapProject[i.Pkg("godep.example.com/app")] = "a"
	project, err := prov.getProject(i.Pkg("godep.example.com/app"))
	require.NoError(err)
	require.Equal("a", project)
}

func TestBBRepo_SetProjectWithError(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://godep.example.com/app?go-get=1").Return([]byte(``), errors.New("err"))

	prov := New(mHttpClient, "bitbucket.example.com", "godep.example.com", "https://bitbucket.example.com")
	project, err := prov.getProject(i.Pkg("godep.example.com/app"))
	require.EqualError(err, "err")
	require.Equal("", project)
}

// func TestGithubRepo_WithDep(t *testing.T) {
// 	t.Parallel()
// 	require := require.New(t)

// 	pkg := i.Pkg("github.com/Masterminds/glide")

// 	mHttpClient := &mocks.IWebClient{}
// 	mHttpClient.On("Get", "Masterminds/glide/master/glide.lock").Return([]byte(`imports:
// - name: pkg1
//   version: hash1`), nil)

// 	prov := New(mHttpClient)
// 	content, err := prov.File(pkg, "glide.lock")
// 	require.NoError(err)
// 	require.True(len(content) > 0)

// 	detector := deps.NewDetector()
// 	detector.AddTool(glide.New())

// 	app := &mocks.IApp{}
// 	app.On("Provider").Return(prov)
// 	app.On("Package").Return(pkg)

// 	appDeps, err := detector.Deps(app)
// 	require.NoError(err)
// 	deps := appDeps.Deps
// 	require.Len(deps, 1, "Expect 1 dependency")
// 	require.Equal(i.Pkg("pkg1"), deps["pkg1"].Package)
// 	require.Equal(i.Hash("hash1"), deps["pkg1"].Hash)
// }

// func TestGithubTests_Ok(t *testing.T) {
// 	t.Parallel()
// 	require := require.New(t)

// 	mHttpClient := &mocks.IWebClient{}
// 	mHttpClient.On("Get", "https://api.github.com/repos/golang/dep/tags").Return([]byte(`[
//   {
//     "name": "v0.1.0",
//     "commit": {
//       "sha": "05c40eba7fa5512c3a161e4e9df6c8fefde75158"
//     }
//   },
//   {
//     "name": "v0.2.0",
//     "commit": {
//       "sha": "6a17782ed25ba4ccd2adf191990cc32e65c3934c"
//     }
//   }
// ]`), nil)
// 	pkg := i.Pkg("github.com/golang/dep")
// 	tagsGetter := New(mHttpClient)

// 	tags, err := tagsGetter.Tags(pkg)
// 	require.NoError(err)
// 	require.Len(tags, 2, "Expect 2 tags")
// 	require.Equal("v0.1.0", tags[0].Version)
// 	require.Equal("v0.2.0", tags[1].Version)
// }

// func TestGithubTags_Error(t *testing.T) {
// 	t.Parallel()
// 	require := require.New(t)

// 	mHttpClient := &mocks.IWebClient{}
// 	mHttpClient.On("Get", "https://api.github.com/repos/golang/dep/tags").Return(nil, errors.New("error"))

// 	pkg := i.Pkg("github.com/golang/dep")
// 	tagsGetter := New(mHttpClient)

// 	tags, err := tagsGetter.Tags(pkg)
// 	require.EqualError(err, "error")
// 	require.Len(tags, 0, "Expect 0 tags")
// }

// func TestGithubTags_BadFile(t *testing.T) {
// 	t.Parallel()
// 	require := require.New(t)

// 	mHttpClient := &mocks.IWebClient{}
// 	mHttpClient.On("Get", "https://api.github.com/repos/golang/dep/tags").Return([]byte(`{`), nil)

// 	pkg := i.Pkg("github.com/golang/dep")
// 	tagsGetter := New(mHttpClient)

// 	tags, err := tagsGetter.Tags(pkg)
// 	require.Error(err)
// 	require.Len(tags, 0, "Expect 0 tags")
// }

// func TestGithubTags_WithToken(t *testing.T) {
// 	t.Parallel()
// 	require := require.New(t)

// 	mHttpClient := &mocks.IWebClient{}
// 	mHttpClient.On("Get", "https://api.github.com/repos/golang/dep/tags").Return([]byte(`{`), nil)

// 	pkg := i.Pkg("github.com/golang/dep")
// 	tagsGetter := New(mHttpClient)

// 	tags, err := tagsGetter.Tags(pkg)
// 	require.Error(err)
// 	require.Len(tags, 0, "Expect 0 tags")
// }

func TestBBRepo_CheckGoGetUrl(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	prov := New(&mocks.IWebClient{}, "bitbucket.example.com", "godep.example.com", "https://bitbucket.example.com")
	require.Equal("godep.example.com", prov.GoGetUrl())
}

// func TestHTTPWrapper_WithoutToken_NoToken(t *testing.T) {
// 	t.Parallel()
// 	require := require.New(t)

// 	mHttpClient := &mocks.IWebClient{}
// 	mHttpClient.On("Get", "tags").Return([]byte(`res`), nil)

// 	client := &HTTPWrapper{
// 		client: mHttpClient,
// 	}
// 	content, err := client.Get("tags")
// 	require.NoError(err)
// 	require.EqualValues("res", content)
// }

// func TestHTTPWrapper_WithToken_AddToken(t *testing.T) {
// 	t.Parallel()
// 	require := require.New(t)

// 	mHttpClient := &mocks.IWebClient{}
// 	mHttpClient.On("Get", "tags?access_token=a").Return([]byte(`res`), nil)

// 	client := &HTTPWrapper{
// 		token:  "a",
// 		client: mHttpClient,
// 	}
// 	content, err := client.Get("tags")
// 	require.NoError(err)
// 	require.EqualValues("res", content)
// }
