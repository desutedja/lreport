package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetStringOrPanic(t *testing.T) {
	t.Run("HaveValue", func(t *testing.T) {
		key := "TestGetStringOrPanic"
		viper.Set(key, "test")
		expected := "test"
		assert.Equal(t, expected, getStringOrPanic(key))
	})
	t.Run("Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			getStringOrPanic("panicValue")
		})
	})
}

func TestGetStringWithDefault(t *testing.T) {
	t.Run("HaveValue", func(t *testing.T) {
		key := "TestGetStringWithDefault"
		viper.Set(key, "test")
		expected := "test"
		assert.Equal(t, expected, getStringWithDefault(key, "default"))
	})
	t.Run("TestGetStringWithDefault", func(t *testing.T) {
		key := "TestGetStringWithDefault2"
		expected := "default"
		assert.Equal(t, expected, getStringWithDefault(key, "default"))
	})
}

func TestGetIntOrPanic(t *testing.T) {
	t.Run("HaveValue", func(t *testing.T) {
		key := "TestGetIntOrPanic"
		viper.Set(key, "1")
		expected := 1
		assert.Equal(t, expected, getIntOrPanic(key))
	})
	t.Run("Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			getIntOrPanic("TestGetIntOrPanic_Panic")
		})
	})
}

func TestGetIntWithDefault(t *testing.T) {
	t.Run("HaveValue", func(t *testing.T) {
		key := "TestGetIntWithDefault"
		viper.Set(key, "123")
		expected := 123
		assert.Equal(t, expected, getIntWithDefault(key, 1))
	})
	t.Run("TestGetIntWithDefault", func(t *testing.T) {
		key := "TestGetIntWithDefault2"
		expected := 1
		assert.Equal(t, expected, getIntWithDefault(key, 1))
	})
}

func TestGetBoolOrPanic(t *testing.T) {
	t.Run("HaveValue", func(t *testing.T) {
		key := "TestGetBoolOrPanic"
		viper.Set(key, "true")
		expected := true
		assert.Equal(t, expected, getBoolOrPanic(key))
	})
	t.Run("Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			getIntOrPanic("TestGetBoolOrPanic_Panic")
		})
	})
}

func TestGetBoolWithDefault(t *testing.T) {
	t.Run("HaveValue", func(t *testing.T) {
		key := "TestGetBoolWithDefault"
		viper.Set(key, "false")
		expected := false
		assert.Equal(t, expected, getBoolWithDefault(key, true))
	})
	t.Run("TestGetIntWithDefault", func(t *testing.T) {
		key := "TestGetIntWithDefault2"
		expected := true
		assert.Equal(t, expected, getBoolWithDefault(key, true))
	})
}
