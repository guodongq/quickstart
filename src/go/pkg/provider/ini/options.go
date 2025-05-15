package ini

func WithIniOptionsEnabled() func(*IniOptions) {
	return func(o *IniOptions) {
		o.Enabled = true
	}
}

func WithIniOptionsIniFile(file string) func(*IniOptions) {
	return func(o *IniOptions) {
		o.IniFile = file
	}
}

func WithIniOptionsSection(section string) func(*IniOptions) {
	return func(o *IniOptions) {
		o.Section = section
	}
}

type IniOptions struct {
	Enabled bool
	IniFile string
	Section string
}

func getDefaultIniOptions() IniOptions {
	return IniOptions{
		Enabled: false,
		IniFile: "config.ini",
		Section: "dev",
	}
}
