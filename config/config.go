package config

// Config is a struct that can be loaded from a JSON file
// or written to a JSON file
type Config struct {
	MinVersion Version  `json:"min_version"` // TODO: implement marshaller/unmarshaller for Version
	Servers    []Server `json:"servers"`     // A list of servers to listen on
}

func (c *Config) Load(filename string) error {
	// read data from file
	// then call json unmarshal
}

func (c *Config) Write(filename string) error {
	// call json marshal
	// then write data to file
}
