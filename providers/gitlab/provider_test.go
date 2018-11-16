package gitlab

import (
	"testing"

	"github.com/stamm/dep_radar"
	"github.com/stretchr/testify/require"
)

func TestCleanID_Ok(t *testing.T) {
	require := require.New(t)
	p := &Provider{
		goGetURL: "gitlab.com",
	}

	var pkgId = "wxcsdb88/go"
	var pkg = dep_radar.Pkg("gitlab.com/" + pkgId)

	require.EqualValues(p.cleanID(pkg), pkgId)
}
