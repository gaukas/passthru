package config

import (
	"fmt"
	"strconv"
	"strings"
        "github.com/gaukas/passthru/internal/logger"
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

const (
	WILL_FIT   uint8 = iota // debug level
	SHOULD_FIT              // info
	MAY_FIT                 // warning
	WONT_FIT                // fatal
)

func (v *Version) CanFitInServer(serverVer *Version) uint8 {
	// for v0, if version is at all different, it won't fit
	if v.Major == 0 && serverVer.Major == 0 {
		if v.Minor != serverVer.Minor || v.Patch != serverVer.Patch {
			return WONT_FIT
		}
		return WILL_FIT
	}

	// v1 and above
	if v.Major != serverVer.Major {
		return WONT_FIT // only same major version can fit!
	}
	if v.Minor > serverVer.Minor { // same major, but minor is higher than server
		return MAY_FIT // should warn user, some features in config may not be available in server
	}
	if v.Minor == serverVer.Minor && v.Patch > serverVer.Patch { // same major/minor, but patch is higher than server
		return SHOULD_FIT // should notifiy user, unintended behavior is possible due to patching
	}
	return WILL_FIT // no problem
}

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
                logger.Errorf("invalid version string: %s", string(data))
		return fmt.Errorf("invalid version string: %s", string(data))
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
                logger.Errorf("invalid version string: %s", string(data))
		return fmt.Errorf("invalid version string: %s", string(data))
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
                logger.Errorf("invalid version string: %s", string(data))
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
