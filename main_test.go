package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetWhitelistEmpty(t *testing.T) {
	got := getWhitelist("")
	want := make([]string, 0)
	assert.EqualValues(t, got, want)
}

func TestGetWhitelistOneElement(t *testing.T) {
	got := getWhitelist("23456ytewrtytr43456")
	want := []string{"23456ytewrtytr43456"}
	assert.EqualValues(t, got, want)
}

func TestGetWhitelistTwoElements(t *testing.T) {
	got := getWhitelist("23456ytewrtytr43456:345tgfdstyuio8765rf")
	want := []string{"23456ytewrtytr43456", "345tgfdstyuio8765rf"}
	assert.EqualValues(t, got, want)
}

func TestIsWhitelistedTrueOneElement(t *testing.T) {
	whitelist = []string{"3456789876543234567"}
	whiteListed := isWhitelisted(3456789876543234567)
	assert.True(t, whiteListed)
}

func TestIsWhitelistedFalseOneElement(t *testing.T) {
	whitelist = []string{"3456789876543234567"}
	whiteListed := isWhitelisted(3456789876543234568)
	assert.False(t, whiteListed)
}

func TestIsWhitelistedTrueMultipleElements(t *testing.T) {
	whitelist = []string{"3456789876543234567", "45678998765465"}
	whiteListed := isWhitelisted(45678998765465)
	assert.True(t, whiteListed)
}

func TestIsWhitelistedFalseMultipleElements(t *testing.T) {
	whitelist = []string{"3456789876543234567", "45678998765465"}
	whiteListed := isWhitelisted(87654334567)
	assert.False(t, whiteListed)
}

func TestGetenvWithValue(t *testing.T) {
	envKey := "TEST_ENV_SETTING"
	defer os.Unsetenv(envKey)
	envValue := "test"
	err := os.Setenv(envKey, envValue)
	if err != nil {
		assert.Error(t, err, "Could not set environment variable")
	}
	variable := getenv(envKey)
	assert.Equal(t, variable, envValue)
}

func TestGetenvNoValue(t *testing.T) {
	envKey := "TEST_ENV_SETTING"
	assert.Panics(t, func() { getenv(envKey) })
}
