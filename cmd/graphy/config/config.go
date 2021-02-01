package config

import (
	"github.com/spf13/viper"
	"strings"
)

var (
	Neo4JHost             = "localhost"
	Neo4JUser             = "neo4j"
	Neo4jPassword         = "password"
	Port                  = 9092
	RequestDrainInSeconds = 5
	PlaygroundEnabled     = false
)

func Init(v *viper.Viper) {
	v.SetDefault("neo4j-host", Neo4JHost)
	v.SetDefault("neo4j-user", Neo4JUser)
	v.SetDefault("neo4j-password", Neo4jPassword)
	v.SetDefault("port", Port)
	v.SetDefault("request-drain-seconds", RequestDrainInSeconds)
	v.SetDefault("playground-enabled", false)

	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.SetEnvPrefix("graphy")
	v.BindEnv("neo4j-host")
	v.BindEnv("neo4j-user")
	v.BindEnv("neo4j-password")
	v.BindEnv("port")
	v.BindEnv("request-drain-seconds")
	v.BindEnv("playground-enabled")

	v.AutomaticEnv()
}

func Load(v *viper.Viper) {
	Neo4JHost = v.GetString("neo4j-host")
	Neo4JUser = v.GetString("neo4j-user")
	Neo4jPassword = v.GetString("neo4j-password")
	Port = v.GetInt("port")
	RequestDrainInSeconds = v.GetInt("request-drain-seconds")
	PlaygroundEnabled = v.GetBool("playground-enabled")
}
