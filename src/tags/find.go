package tags

import (
	"fmt"

	i "github.com/stamm/dep_radar/interfaces"
)

func FindByHash(hash i.Hash, tags []i.Tag) (i.Tag, error) {
	for _, tag := range tags {
		if tag.Hash == hash {
			return tag, nil
		}
	}
	return i.Tag{}, fmt.Errorf("can't find hash %s", hash)
}

func FindByVersion(version string, tags []i.Tag) (i.Tag, error) {
	for _, tag := range tags {
		if tag.Version == version {
			return tag, nil
		}
	}
	return i.Tag{}, fmt.Errorf("can't find version %s", version)
}
