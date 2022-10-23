package config

import (
	"fmt"
	"strconv"
	"strings"
)

// Version is a struct that represents the version of the config file
// v1.2.3 would be represented as:
//
//	Version{
//		Major: 1,
//		Minor: 2,
//		Patch: 3,
//	}
type Version struct {
	Major int
	Minor int
	Patch int
}

// Will have to implement the custom unmarshaller and marshaller for Version
// due to type conflict. (JSON: string, Go: struct with 3 ints)

// Read: https://mariadesouza.com/2017/09/07/custom-unmarshal-json-in-golang/

func (v *Version) UnmarshalJSON(data []byte) error {
	// v1.2.3 = v$Major.$Minor.$Patch
	strver := string(data)

	// remove the "v" prefix and \"
	strver = strings.TrimPrefix(strver, "\"v")
	strver = strings.TrimSuffix(strver, "\"")

	// split the string into 3 parts
	// major, minor, patch
	parts := strings.Split(strver, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid version string: %s", string(data))
	}

	// convert the parts into ints
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid version string: %s", string(data))
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid version string: %s", string(data))
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid version string: %s", string(data))
	}

	v.Major = major
	v.Minor = minor
	v.Patch = patch

	return nil
}

func (v *Version) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"v%d.%d.%d\"", v.Major, v.Minor, v.Patch)), nil
}
