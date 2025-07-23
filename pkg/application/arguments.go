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
	Controller  string `mapstructure:"controller,omitempty"`
	Model       string `mapstructure:"model,omitempty"`
	Application string `mapstructure:"application,omitempty"`
	Key         string `mapstructure:"key,omitempty"`
	Value       string `mapstructure:"value,omitempty"`
}

type AddModelToolArgs struct {
	Controller  string `mapstructure:"controller,omitempty"`
	Model       string `mapstructure:"model,omitempty"`
	Owner       string `mapstructure:"owner,omitempty"`
	Config      string `mapstructure:"config,omitempty"`
	Credential  string `mapstructure:"credential,omitempty"`
	CloudRegion string `mapstructure:"cloud-region,omitempty"`
	NoSwitch    bool   `mapstructure:"no-switch,omitempty"`
}
