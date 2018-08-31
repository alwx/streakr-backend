package root

import (
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

func ReadConfig() {
	viper.SetDefault("environment", "development")
	viper.SetEnvPrefix("streakr")
	viper.BindEnv("environment")

	viper.SetConfigName(viper.GetString("environment"))
	viper.AddConfigPath("/etc/streakr-go/configs")
	viper.AddConfigPath("./configs")

	log.Println(viper.GetString("environment"))

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while loading config: %s", err)
	}
}
