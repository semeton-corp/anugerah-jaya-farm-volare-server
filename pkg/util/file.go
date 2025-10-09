package util

import (
	"fmt"

	"github.com/spf13/viper"
)

func ConstructURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", viper.GetString("s3.endpoint"), viper.GetString("s3.bucket"), key)
}
