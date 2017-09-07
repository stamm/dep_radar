package bitbucketprivate

import (
	"errors"
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/interfaces/mocks"
	"github.com/stamm/dep_radar/src/deps"
	"github.com/stamm/dep_radar/src/deps/glide"
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

func TestBBRepo_WithDep(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	pkg := i.Pkg("godep.example.com/app")

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://bitbucket.example.com/projects/go/repos/app/raw/glide.lock").Return([]byte(`imports:
- name: pkg1
  version: hash1`), nil)

	prov := New(mHttpClient, "bitbucket.example.com", "godep.example.com", "https://bitbucket.example.com")
	prov.mapProject[pkg] = "go"
	content, err := prov.File(pkg, "glide.lock")
	require.NoError(err)
	require.True(len(content) > 0)

	detector := deps.NewDetector()
	detector.AddTool(glide.New())

	app := &mocks.IApp{}
	app.On("Provider").Return(prov)
	app.On("Package").Return(pkg)

	appDeps, err := detector.Deps(app)
	require.NoError(err)
	deps := appDeps.Deps
	require.Len(deps, 1, "Expect 1 dependency")
	require.Equal(i.Pkg("pkg1"), deps["pkg1"].Package)
	require.Equal(i.Hash("hash1"), deps["pkg1"].Hash)
}

// Tags

func TestBBTags_Ok(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://bitbucket.example.com/rest/api/1.0/projects/go/repos/app/tags?start=0").Return([]byte(`
 {"size":25,"limit":25,"isLastPage":true,"values":[
 {"id":"refs/tags/v0.1.0","displayId":"10.4.0","type":"TAG","latestCommit":"5d6de801e6bb31eb1086adf3604570693580f141","latestChangeset":"5d6de801e6bb31eb1086adf3604570693580f141","hash":null},
 {"id":"refs/tags/v0.2.0","displayId":"10.4.0","type":"TAG","latestCommit":"5d6de801e6bb31eb1086adf3604570693580f141","latestChangeset":"5d6de801e6bb31eb1086adf3604570693580f141","hash":null}
 ],"start":0} `), nil)
	pkg := i.Pkg("godep.example.com/app")
	tagsGetter := New(mHttpClient, "bitbucket.example.com", "godep.example.com", "https://bitbucket.example.com")
	tagsGetter.mapProject[pkg] = "go"

	tags, err := tagsGetter.Tags(pkg)
	require.NoError(err)
	require.Len(tags, 2, "Expect 2 tags")
	require.Equal("v0.1.0", tags[0].Version)
	require.Equal("v0.2.0", tags[1].Version)
}

func TestBBTags_TwoPages(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://bitbucket.example.com/rest/api/1.0/projects/go/repos/app/tags?start=0").Return([]byte(`
 {"size":1,"limit":1,"isLastPage":false,"values":[
 {"id":"refs/tags/v0.1.0","displayId":"10.4.0","type":"TAG","latestCommit":"5d6de801e6bb31eb1086adf3604570693580f141","latestChangeset":"5d6de801e6bb31eb1086adf3604570693580f141","hash":null}
 ],"start":0,"nextPageStart":1} `), nil)
	mHttpClient.On("Get", "https://bitbucket.example.com/rest/api/1.0/projects/go/repos/app/tags?start=1").Return([]byte(`
 {"size":1,"limit":1,"isLastPage":true,"values":[
 {"id":"refs/tags/v0.2.0","displayId":"10.4.0","type":"TAG","latestCommit":"5d6de801e6bb31eb1086adf3604570693580f141","latestChangeset":"5d6de801e6bb31eb1086adf3604570693580f141","hash":null}
 ],"start":1} `), nil)
	pkg := i.Pkg("godep.example.com/app")
	tagsGetter := New(mHttpClient, "bitbucket.example.com", "godep.example.com", "https://bitbucket.example.com")
	tagsGetter.mapProject[pkg] = "go"

	tags, err := tagsGetter.Tags(pkg)
	require.NoError(err)
	require.Len(tags, 2, "Expect 2 tags")
	require.Equal("v0.1.0", tags[0].Version)
	require.Equal("v0.2.0", tags[1].Version)
}

func TestGithubTags_Error(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://bitbucket.example.com/rest/api/1.0/projects/go/repos/app/tags?start=0").Return(nil, errors.New("error"))

	pkg := i.Pkg("godep.example.com/app")
	tagsGetter := New(mHttpClient, "bitbucket.example.com", "godep.example.com", "https://bitbucket.example.com")
	tagsGetter.mapProject[pkg] = "go"

	tags, err := tagsGetter.Tags(pkg)
	require.EqualError(err, "Error on getting tags: error")
	require.Len(tags, 0, "Expect 0 tags")
}

func TestBBTags_BadFile(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://bitbucket.example.com/rest/api/1.0/projects/go/repos/app/tags?start=0").Return([]byte(`{`), nil)

	pkg := i.Pkg("godep.example.com/app")
	tagsGetter := New(mHttpClient, "bitbucket.example.com", "godep.example.com", "https://bitbucket.example.com")
	tagsGetter.mapProject[pkg] = "go"

	tags, err := tagsGetter.Tags(pkg)
	require.Error(err)
	require.Len(tags, 0, "Expect 0 tags")
}

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
