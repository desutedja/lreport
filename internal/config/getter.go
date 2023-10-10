package config

import (
	"log"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

var noConfigErr = "config with key %s is not set"

func getStringWithDefault(key string, defaultValue string) string {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getStringOrPanic(key string) string {
	value := viper.GetString(key)
	if value == "" {
		log.Panicf(noConfigErr, key)
	}
	return value
}

func getIntWithDefault(key string, defaultValue int) int {
	val := viper.GetString(key)

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return intVal
}

func getIntOrPanic(key string) int {
	val := viper.GetString(key)

	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Panicf("%s is not integer value", key)
	}

	return intVal
}

func getBoolWithDefault(key string, defaultValue bool) bool {
	val := viper.GetString(key)

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}

	return boolVal
}

func getBoolOrPanic(key string) bool {
	val := viper.GetString(key)

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		log.Panicf("%s is not boolean value", key)
	}

	return boolVal
}

func getDurationOrPanic(key string) time.Duration {
	val := viper.GetString(key)

	durationVal, err := time.ParseDuration(val)
	if err != nil {
		log.Panicf("%s is not valid duration", key)
	}

	return durationVal
}
