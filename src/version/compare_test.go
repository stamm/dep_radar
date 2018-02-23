package version_test

import (
	"fmt"
	"testing"

	"github.com/stamm/dep_radar/src/version"
	"github.com/stretchr/testify/require"
)

type dataType struct {
	Recommended string
	Actual      string
	Result      bool
}

func TestWithEmpty(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	ok, err := version.Compare("", "1.0.0")
	require.NoError(err)
	require.False(ok)

	ok, err = version.Compare("1.0.0", "")
	require.NoError(err)
	require.False(ok)
}

func TestWithError(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	ok, err := version.Compare("erunda", "1.0.0")
	require.Error(err)
	require.False(ok)

	ok, err = version.Compare("1.0.0", "erunda")
	require.Error(err)
	require.False(ok)
}

func TestCompareVersion(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	data := []dataType{
		{"0.5.0", "1.0.0", false},
		{"^0.5.0", "1.0.0", false},
		{"^0.5.0", "0.6.0", true},
		{"~0.5.0", "1.0.0", false},
		{"~0.5.0", "0.6.0", false},
		{"~0.5.0", "0.5.1", true},
		{">=0.5.0", "1.0.0", true},
		{">=0.5.0", "0.4.9", false},
		{">0.5.0", "0.5.0", false},
		{">0.5.0", "0.5.1", true},
		{"<0.5.0", "0.5.0", false},
		{"<0.5.0", "1.0.0", false},
		{"<=0.5.0", "0.4.0", true},

		{">=0.5.0", "v0.4.0", false},
		{">=0.5.0", "v0.5.0", true},

		{">=v0.5.0", "0.4.0", false},
		{">=v0.5.0", "0.5.0", true},

		{">0.5.1|>=1.0.0", "0.5.2", true},
		{">0.5.1|>=1.0.0", "0.5.1", false},
		{">0.5.1|>=1.0.0", "0.5.0", false},

		{"^0.5.1|>=1.2.0", "1.2.0", true},
		{"^0.5.1|>=1.2.0", "1.1.0", false},
	}
	for _, line := range data {
		ok, err := version.Compare(line.Recommended, line.Actual)
		require.NoError(err, fmt.Sprintf("Not expect error %s on Compare(%s, %s)", err, line.Recommended, line.Actual))
		require.Equal(line.Result, ok, fmt.Sprintf("Wrong value %t instead of %t on Compare(%s, %s)", ok, line.Result, line.Recommended, line.Actual))
	}
}
