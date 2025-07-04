package config

type app struct {
	Application application `mapstructure:"application"`
	Logger      logger      `mapstructure:"logger"`
	User        user        `mapstructure:"user"`
	Kafka       kafka       `mapstructure:"kafka"`
	Email       email       `mapstructure:"email"`
}

// Application section
type application struct {
	Name string `mapstructure:"name"`
}

// Logger section
type logger struct {
	Level    string `mapstructure:"level"`
	Encoding struct {
		Method string `mapstructure:"method"`
		Caller bool   `mapstructure:"caller"`
	} `mapstructure:"encoding"`
	Path    string `mapstructure:"path"`
	ErrPath string `mapstructure:"errPath"`
}

// Kafka section
type kafka struct {
	Brokers []string `mapstructure:"brokers"`
	Topics  struct {
		Activation    string `mapstructure:"activation"`
		PasswordReset string `mapstructure:"passwordReset"`
	} `mapstructure:"topics"`
	ConsumerGroupName string `mapstructure:"consumer_group_name"`
}

type user struct {
	Activation struct {
		TemplatePath string `mapstructure:"templatePath"`
	} `mapstructure:"activation"`
	PasswordReset struct {
		TemplatePath string `mapstructure:"templatePath"`
	} `mapstructure:"passwordReset"`
}

type email struct {
	SMTP struct {
		From     string `mapstructure:"from"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
	} `mapstructure:"smtp"`
}
