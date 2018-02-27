package github

import (
	"context"
	"errors"
	"testing"

	"github.com/stamm/dep_radar/src/deps"
	"github.com/stamm/dep_radar/src/deps/glide"
	i "github.com/stamm/dep_radar/src/interfaces"
	"github.com/stamm/dep_radar/src/interfaces/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGithubRepo_GetUrl(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	prov := New(&mocks.IWebClient{})
	url := prov.makeURL(i.Pkg("github.com/Masterminds/glide"), "master", "glide.lock")
	require.Equal("https://raw.githubusercontent.com/Masterminds/glide/master/glide.lock", url)
}

func TestGithubRepo_GetUrlWithBranch(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	prov := New(&mocks.IWebClient{})
	url := prov.makeURL(i.Pkg("github.com/Masterminds/glide"), "dev", "glide.lock")
	require.Equal("https://raw.githubusercontent.com/Masterminds/glide/dev/glide.lock", url)
}

func TestGithubRepo_GetUrlSubpackage(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	app := New(&mocks.IWebClient{})
	url := app.makeURL(i.Pkg("github.com/stretchr/testify/require"), "master", "doc.go")
	require.Equal("https://raw.githubusercontent.com/stretchr/testify/master/require/doc.go", url)
}

func TestGithubRepo_GetUrlWithSlash(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	app := New(&mocks.IWebClient{})
	url := app.makeURL(i.Pkg("github.com/stretchr/testify/require/"), "master", "doc.go")
	require.Equal("https://raw.githubusercontent.com/stretchr/testify/master/require/doc.go", url)
}

func TestGithubRepo_WithDep(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	pkg := i.Pkg("github.com/Masterminds/glide")

	client := &mocks.IWebClient{}
	client.On("Get", mock.Anything, "https://raw.githubusercontent.com/Masterminds/glide/master/glide.lock").Return([]byte(`imports:
- name: pkg1
  version: hash1`), nil)

	prov := New(client)
	content, err := prov.File(context.Background(), pkg, "master", "glide.lock")
	require.NoError(err)
	require.True(len(content) > 0)

	detector := deps.NewDetector()
	detector.AddTool(glide.New())

	app := &mocks.IApp{}
	app.On("Provider").Return(prov)
	app.On("Package").Return(pkg)
	app.On("Branch").Return("master")

	appDeps, err := detector.Deps(context.Background(), app)
	require.NoError(err)
	deps := appDeps.Deps
	require.Len(deps, 1, "Expect 1 dependency")
	require.Equal(i.Pkg("pkg1"), deps["pkg1"].Package)
	require.Equal(i.Hash("hash1"), deps["pkg1"].Hash)
}

func TestGithubRepo_CheckGoGetURL(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	require.Equal("github.com", New(&mocks.IWebClient{}).GoGetURL())
}

// Tags

func TestGithubTags_Ok(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &mocks.IWebClient{}
	client.On("Get", mock.Anything, "https://api.github.com/repos/golang/dep/tags?per_page=100").Return([]byte(`[
  {
    "name": "v0.1.0",
    "commit": {
      "sha": "05c40eba7fa5512c3a161e4e9df6c8fefde75158"
    }
  },
  {
    "name": "v0.2.0",
    "commit": {
      "sha": "6a17782ed25ba4ccd2adf191990cc32e65c3934c"
    }
  }
]`), nil)
	pkg := i.Pkg("github.com/golang/dep")
	tagsGetter := New(client)

	tags, err := tagsGetter.Tags(context.Background(), pkg)
	require.NoError(err)
	require.Len(tags, 2, "Expect 2 tags")
	require.Equal("v0.1.0", tags[0].Version)
	require.Equal("v0.2.0", tags[1].Version)
}

func TestGithubTags_Error(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &mocks.IWebClient{}
	client.On("Get", mock.Anything, "https://api.github.com/repos/golang/dep/tags?per_page=100").Return(nil, errors.New("error"))

	pkg := i.Pkg("github.com/golang/dep")
	tagsGetter := New(client)

	tags, err := tagsGetter.Tags(context.Background(), pkg)
	require.EqualError(err, "error")
	require.Len(tags, 0, "Expect 0 tags")
}

func TestGithubTags_BadFile(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &mocks.IWebClient{}
	client.On("Get", mock.Anything, "https://api.github.com/repos/golang/dep/tags?per_page=100").Return([]byte(`{`), nil)

	pkg := i.Pkg("github.com/golang/dep")
	tagsGetter := New(client)

	tags, err := tagsGetter.Tags(context.Background(), pkg)
	require.Error(err)
	require.Len(tags, 0, "Expect 0 tags")
}

// HTTPWrapper

func TestHTTPWrapper_WithoutToken_NoToken(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mClient := &mocks.IWebClient{}
	mClient.On("Get", mock.Anything, "tags").Return([]byte(`res`), nil)

	client := &HTTPWrapper{
		client: mClient,
	}
	content, err := client.Get(context.Background(), "tags")
	require.NoError(err)
	require.EqualValues("res", content)
}

func TestHTTPWrapper_WithToken_AddToken(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mClient := &mocks.IWebClient{}
	mClient.On("Get", mock.Anything, "tags?access_token=a").Return([]byte(`res`), nil)

	client := &HTTPWrapper{
		token:  "a",
		client: mClient,
	}
	content, err := client.Get(context.Background(), "tags")
	require.NoError(err)
	require.EqualValues("res", content)
}

func TestHTTPWrapper_GetURL_WithToken_ExpectAdd(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &HTTPWrapper{
		token: "a",
	}
	url, err := client.getURL("tags?test=1")
	require.NoError(err)
	require.EqualValues("tags?test=1&access_token=a", url)
}

func TestHTTPWrapper_GetURL_WrongUrl(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &HTTPWrapper{
		token: "a",
	}
	url, err := client.getURL("cache_object:foo")
	require.Error(err)
	require.EqualValues("", url)
}

func TestHTTPWrapper_Get_WrongURL(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := &HTTPWrapper{
		token: "a",
	}
	content, err := client.Get(context.Background(), "cache_object:foo")
	require.Error(err)
	require.EqualValues("", content)
}

func TestHTTPWrapper_New(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := NewHTTPWrapper("a", 2)
	url, err := client.getURL("tags")
	require.NoError(err)
	require.EqualValues("tags?access_token=a", url)
}
