package config

// Flags are applied after the configuration file is loaded.
// They are pointers to represent optional types, to tell whether they have been set.
type Flags struct {
	Destination *string
	Unpublished *bool
	Drafts      *bool
	Future      *bool
}

// ApplyFlags overwrites the configuration with values from flags.
func (c *Config) ApplyFlags(flags Flags) {
	if flags.Destination != nil {
		c.Destination = *flags.Destination
	}
	if flags.Drafts != nil {
		c.Drafts = *flags.Drafts
	}
	if flags.Future != nil {
		c.Future = *flags.Future
	}
	if flags.Unpublished != nil {
		c.Unpublished = *flags.Unpublished
	}
}
