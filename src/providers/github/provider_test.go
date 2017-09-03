package github

import (
	"errors"
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/interfaces/mocks"
	"github.com/stamm/dep_radar/src/deps/glide"
	"github.com/stretchr/testify/assert"
)

func TestGithubRepo_GetUrl(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	app, err := New("github.com/Masterminds/glide", &mocks.IWebClient{})
	assert.NoError(err)
	url := app.makeURL("glide.lock")
	assert.Equal(url, "Masterminds/glide/master/glide.lock")
}

func TestGithubRepo_WrongPackage(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	_, err := New("gopkg.in/yaml.v2", &mocks.IWebClient{})
	assert.EqualError(err, "package gopkg.in/yaml.v2 is not for github")
}

func TestGithubRepo_GetUrlSubpackage(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	app, err := New("github.com/stretchr/testify/assert", &mocks.IWebClient{})
	assert.NoError(err)
	url := app.makeURL("doc.go")
	assert.Equal(url, "stretchr/testify/master/assert/doc.go")
}

func TestGithubRepo_GetUrlWithSlash(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	app, err := New("github.com/stretchr/testify/assert/", &mocks.IWebClient{})
	assert.NoError(err)
	url := app.makeURL("doc.go")
	assert.Equal(url, "stretchr/testify/master/assert/doc.go")
}

func TestGithubRepo(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "Masterminds/glide/master/glide.lock").Return([]byte(`imports:
- name: pkg1
  version: hash1`), nil)

	app, err := New("github.com/Masterminds/glide", mHttpClient)
	assert.NoError(err)
	assert.Equal(app.Package(), i.Pkg("github.com/Masterminds/glide"))
	content, err := app.File("glide.lock")
	assert.NoError(err)
	assert.True(len(content) > 0)

	appDeps, err := glide.Tool(app)
	assert.NoError(err)
	deps := appDeps.Deps
	assert.Equal(len(deps), 1, "Expect 1 dependency")
	assert.Equal(deps["pkg1"].Package, "pkg1")
	assert.Equal(deps["pkg1"].Hash, i.Hash("hash1"))
}

func TestGithubTests_Ok(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
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
	tagsGetter, err := New("github.com/golang/dep", mHttpClient)
	assert.NoError(err)

	tags, err := tagsGetter.Tags()
	assert.NoError(err)
	assert.Equal(len(tags), 2, "Expect 2 tags")
	assert.Equal(tags[0].Version, "v0.1.0")
	assert.Equal(tags[1].Version, "v0.2.0")
	// assert.Equal(deps[0].Package, "pkg1")
	// assert.Equal(deps[0].Hash, i.Hash("hash1"))
}

func TestGithubTags_Error(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://api.github.com/repos/golang/dep/tags").Return(nil, errors.New("error"))
	tagsGetter, err := New("github.com/golang/dep", mHttpClient)
	assert.NoError(err)

	tags, err := tagsGetter.Tags()
	assert.EqualError(err, "error")
	assert.Equal(len(tags), 0, "Expect 0 tags")
}

func TestGithubTags_BadFile(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://api.github.com/repos/golang/dep/tags").Return([]byte(`{`), nil)
	tagsGetter, err := New("github.com/golang/dep", mHttpClient)
	assert.NoError(err)

	tags, err := tagsGetter.Tags()
	assert.Error(err)
	assert.Equal(len(tags), 0, "Expect 0 tags")
}

func TestGithubTags_WithToken(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	mHttpClient := &mocks.IWebClient{}
	mHttpClient.On("Get", "https://api.github.com/repos/golang/dep/tags").Return([]byte(`{`), nil)

	tagsGetter, err := New("github.com/golang/dep", mHttpClient)
	assert.NoError(err)

	tags, err := tagsGetter.Tags()
	assert.Error(err)
	assert.Equal(len(tags), 0, "Expect 0 tags")
}
