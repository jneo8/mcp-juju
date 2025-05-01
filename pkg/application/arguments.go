package application

type ListModelToolArgs struct {
	Controller string `mapstructure:"controller,omitempty"`
}

type GetStatusToolArgs struct {
	Controller     string `mapstructure:"controller,omitempty"`
	Model          string `mapstructure:"model,omitempty"`
	IncludeStorage bool   `mapstructure:"include-storage,omitempty"`
}

type GetApplicationConfigToolArgs struct {
	Controller  string `mapstructure:"controller,omitempty"`
	Model       string `mapstructure:"model,omitempty"`
	Application string `mapstructure:"application,omitempty"`
}

type SetApplicationConfigToolArgs struct {
	Controller  string            `mapstructure:"controller,omitempty"`
	Model       string            `mapstructure:"model,omitempty"`
	Application string            `mapstructure:"application,omitempty"`
	Settings    map[string]string `mapstructure:"settings,omitempty"`
}
