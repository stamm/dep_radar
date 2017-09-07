package github

import (
	"errors"
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/interfaces/mocks"
	"github.com/stamm/dep_radar/src/deps"
	"github.com/stamm/dep_radar/src/deps/glide"
	"github.com/stretchr/testify/require"
)

func TestGithubRepo_GetUrl(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	prov := New(&mocks.IWebClient{})
	url := prov.makeURL(i.Pkg("github.com/Masterminds/glide"), "glide.lock")
	require.Equal("Masterminds/glide/master/glide.lock", url)
}

func TestGithubRepo_GetUrlSubpackage(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	app := New(&mocks.IWebClient{})
	url := app.makeURL(i.Pkg("github.com/stretchr/testify/require"), "doc.go")
	require.Equal("stretchr/testify/master/require/doc.go", url)
}

func TestGithubRepo_GetUrlWithSlash(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	app := New(&mocks.IWebClient{})
	url := app.makeURL(i.Pkg("github.com/stretchr/testify/require/"), "doc.go")
	require.Equal("stretchr/testify/master/require/doc.go", url)
}

func TestGithubRepo_WithDep(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	pkg := i.Pkg("github.com/Masterminds/glide")

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "Masterminds/glide/master/glide.lock").Return([]byte(`imports:
- name: pkg1
  version: hash1`), nil)

	prov := New(mHttpClient)
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

func TestGithubRepo_CheckGoGetUrl(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	require.Equal("github.com", New(&mocks.IWebClient{}).GoGetUrl())
}

// Tags

func TestGithubTags_Ok(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://api.github.com/repos/golang/dep/tags").Return([]byte(`[
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
	tagsGetter := New(mHttpClient)

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
	mHttpClient.On("Get", "https://api.github.com/repos/golang/dep/tags").Return(nil, errors.New("error"))

	pkg := i.Pkg("github.com/golang/dep")
	tagsGetter := New(mHttpClient)

	tags, err := tagsGetter.Tags(pkg)
	require.EqualError(err, "error")
	require.Len(tags, 0, "Expect 0 tags")
}

func TestGithubTags_BadFile(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://api.github.com/repos/golang/dep/tags").Return([]byte(`{`), nil)

	pkg := i.Pkg("github.com/golang/dep")
	tagsGetter := New(mHttpClient)

	tags, err := tagsGetter.Tags(pkg)
	require.Error(err)
	require.Len(tags, 0, "Expect 0 tags")
}

// HTTPWrapper

func TestHTTPWrapper_WithoutToken_NoToken(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "tags").Return([]byte(`res`), nil)

	client := &HTTPWrapper{
		client: mHttpClient,
	}
	content, err := client.Get("tags")
	require.NoError(err)
	require.EqualValues("res", content)
}

func TestHTTPWrapper_WithToken_AddToken(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "tags?access_token=a").Return([]byte(`res`), nil)

	client := &HTTPWrapper{
		token:  "a",
		client: mHttpClient,
	}
	content, err := client.Get("tags")
	require.NoError(err)
	require.EqualValues("res", content)
}
