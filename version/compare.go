package version

import (
	"log"
	"strings"

	"github.com/Masterminds/semver"
)

// Compare versions
func Compare(recomends, actual string) (bool, error) {
	if recomends == "" || actual == "" {
		return false, nil
	}
	for _, recommended := range strings.Split(recomends, "|") {
		c, err := semver.NewConstraint(recommended)
		if err != nil {
			log.Printf("err (%s, %s) = %+v\n", recommended, actual, err)
			return false, err
		}
		v, err := semver.NewVersion(actual)
		if err != nil {
			log.Printf("err (%s, %s) = %+v\n", recommended, actual, err)
			return false, err
		}
		if c.Check(v) {
			return true, nil
		}
	}
	return false, nil
}
