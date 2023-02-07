/*
 * Copyright 2023 Asim Ihsan
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package identifier

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	operatorRegex = regexp.MustCompile(`([><=]{1,2})\s*(.*)`)
)

type RequirementType int

const (
	Exact RequirementType = iota
	Caret
	Tilde
	SingleConditionEqual
	SingleConditionGreaterThan
	SingleConditionLessThan
	SingleConditionGreaterThanOrEqual
	SingleConditionLessThanOrEqual
)

type Requirement struct {
	Type       RequirementType
	Version    SemverVersion
	MaxVersion *SemverVersion
}

type SemverVersion struct {
	Major int
	Minor *int
	Patch *int
}

func CompareSemverVersions(a, b SemverVersion) int {
	if a.Major > b.Major {
		return 1
	}
	if a.Major < b.Major {
		return -1
	}
	if a.Minor != nil && b.Minor != nil {
		if *a.Minor > *b.Minor {
			return 1
		}
		if *a.Minor < *b.Minor {
			return -1
		}
	} else if a.Minor != nil {
		return 1
	} else if b.Minor != nil {
		return -1
	}
	if a.Patch != nil && b.Patch != nil {
		if *a.Patch > *b.Patch {
			return 1
		}
		if *a.Patch < *b.Patch {
			return -1
		}
	} else if a.Patch != nil {
		return 1
	} else if b.Patch != nil {
		return -1
	}
	return 0
}

func NewRequirement(s string) (*Requirement, error) {
	if strings.HasPrefix(s, "^") {
		version, err := ParseVersion(s[1:])
		if err != nil {
			return nil, err
		}
		return &Requirement{
			Type:    Caret,
			Version: *version,
		}, nil
	}
	if strings.HasPrefix(s, "~") {
		version, err := ParseVersion(s[1:])
		if err != nil {
			return nil, err
		}
		return &Requirement{
			Type:    Tilde,
			Version: *version,
		}, nil
	}

	conditionOperatorToType := map[string]RequirementType{
		"==": SingleConditionEqual,
		">":  SingleConditionGreaterThan,
		"<":  SingleConditionLessThan,
		">=": SingleConditionGreaterThanOrEqual,
		"<=": SingleConditionLessThanOrEqual,
	}

	matches := operatorRegex.FindStringSubmatch(s)
	if len(matches) == 3 {
		operator := matches[1]
		version := matches[2]

		requirementType, ok := conditionOperatorToType[operator]
		if ok {
			version, err := ParseVersion(version)
			if err != nil {
				return nil, err
			}
			return &Requirement{
				Type:    requirementType,
				Version: *version,
			}, nil
		}
	}

	version, err := ParseVersion(s)
	if err != nil {
		return nil, err
	}

	return &Requirement{
		Type:    Exact,
		Version: *version,
	}, nil
}

func mustParseVersion(s string) *SemverVersion {
	v, err := ParseVersion(s)
	if err != nil {
		panic(err)
	}
	return v
}

func ParseVersion(s string) (*SemverVersion, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "v")

	// Split into major.minor.patch
	parts := strings.SplitN(s, ".", 3)

	if len(parts) > 3 {
		return nil, errors.New("invalid version, too many parts")
	}

	var major int
	var minor int
	var isMinorSet = false
	var patch int
	var isPatchSet = false
	var err error

	if len(parts) >= 1 {
		major, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
	}

	if len(parts) >= 2 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		isMinorSet = true
	}

	if len(parts) == 3 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			return nil, err
		}
		isPatchSet = true
	}

	if isMinorSet && isPatchSet {
		return &SemverVersion{
			Major: major,
			Minor: &minor,
			Patch: &patch,
		}, nil
	} else if isMinorSet {
		return &SemverVersion{
			Major: major,
			Minor: &minor,
		}, nil
	} else {
		return &SemverVersion{
			Major: major,
		}, nil
	}
}

// Satisfies returns true if the version matches the semver. version is the version of the
// program, and requirement is a semver requirement. The semver requirement is a string that
// follows the conventions in https://doc.rust-lang.org/cargo/reference/specifying-dependencies.html.
//
// Examples (version is on the left, requirement is on the right):
// - 1.2.3 matches 1.2.3
// - 1.2.3 matches ^1.2.3
// - 1.2.3 does not match 1.2
// - 1.2.3 matches ~1.2.3
// - 1.2.3 matches ~1.2
// - 1.2.3 matches ~1
// - 1.2.3 does not match ~2
func Satisfies(version string, requirement string) bool {
	req, err := NewRequirement(requirement)
	if err != nil {
		return false
	}
	v, err := ParseVersion(version)
	if err != nil {
		return false
	}

	switch req.Type {
	case Exact, Caret:
		return CompareSemverVersions(*v, req.Version) == 0

	case Tilde:
		// If req only has major version, then major versions must match.
		if req.Version.Minor == nil && req.Version.Patch == nil {
			return req.Version.Major == v.Major
		}

		// If req only has major and minor versions, then major and minor versions must match.
		if req.Version.Patch == nil {
			return req.Version.Major == v.Major && *req.Version.Minor == *v.Minor
		}

		// If req has all of major, minor, and patch, then the version must have the same
		// major and minor versions, and the patch version must be greater than or equal to
		// the patch version in the requirement.
		return req.Version.Major == v.Major &&
			*req.Version.Minor == *v.Minor &&
			CompareSemverVersions(*v, req.Version) >= 0

	case SingleConditionEqual:
		return CompareSemverVersions(*v, req.Version) == 0
	case SingleConditionGreaterThan:
		return CompareSemverVersions(*v, req.Version) > 0
	case SingleConditionLessThan:
		return CompareSemverVersions(*v, req.Version) < 0
	case SingleConditionGreaterThanOrEqual:
		return CompareSemverVersions(*v, req.Version) >= 0
	case SingleConditionLessThanOrEqual:
		return CompareSemverVersions(*v, req.Version) <= 0
	}

	return false
}
