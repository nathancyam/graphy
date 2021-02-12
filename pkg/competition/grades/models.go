package grades

type Grade struct {
	ID  string `mapstructure:"id"`
	Age string `mapstructure:"age"`
	Day string `mapstructure:"day"`
}
