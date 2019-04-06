package util

import (
	"github.com/Masterminds/semver"
)

func CheckVersionConstraint(constraint, version string) (bool, error) {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return false, err
	}
	v, err := semver.NewVersion(version)
	if err != nil {
		return false, err
	}
	return c.Check(v), nil
}
