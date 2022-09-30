package config

// Version is a struct that represents the version of the config file
// v1.2.3 would be represented as:
// Version{
// 	Major: 1,
// 	Minor: 2,
// 	Patch: 3,
// }
type Version struct {
	Major int
	Minor int
	Patch int
}

// Will have to implement the custom unmarshaller and marshaller for Version
// due to type conflict. (JSON: string, Go: struct with 3 ints)

// Read: https://mariadesouza.com/2017/09/07/custom-unmarshal-json-in-golang/

func (v *Version) UnmarshalJSON(data []byte) error {
	return nil
}

func (v *Version) MarshalJSON() ([]byte, error) {
	return nil, nil
}
