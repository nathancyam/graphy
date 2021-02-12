package rounds

type Round struct {
	ID              string `mapstructure:"id" json:"id"`
	Name            string `mapstructure:"name" json:"name"`
	SequenceNo      int    `mapstructure:"sequenceNo" json:"sequenceNo"`
	RoutingCode     string `mapstructure:"routingCode" json:"routingCode"`
	ProvisionalDate string `mapstructure:"provisionalDate" json:"provisionalDate"`
}
