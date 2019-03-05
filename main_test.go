package main

import (
	"gotest.tools/assert"
	"testing"
)

func TestGetWhitelistEmpty(t *testing.T) {
	got := getWhitelist("")
	want := make([]string, 0)
	assert.DeepEqual(t, got, want)
}

func TestGetWhitelistOneElement(t *testing.T) {
	got := getWhitelist("23456ytewrtytr43456")
	want := []string{"23456ytewrtytr43456"}
	assert.DeepEqual(t, got, want)
}

func TestGetWhitelistTwoElements(t *testing.T) {
	got := getWhitelist("23456ytewrtytr43456:345tgfdstyuio8765rf")
	want := []string{"23456ytewrtytr43456", "345tgfdstyuio8765rf"}
	assert.DeepEqual(t, got, want)
}

func TestIsWhitelistedTrueOneElement(t *testing.T) {
	whitelist = []string{"3456789876543234567"}
	whiteListed := isWhitelisted(3456789876543234567)
	assert.Assert(t, whiteListed)
}

func TestIsWhitelistedFalseOneElement(t *testing.T) {
	whitelist = []string{"3456789876543234567"}
	whiteListed := isWhitelisted(3456789876543234568)
	assert.Assert(t, !whiteListed)
}

func TestIsWhitelistedTrueMultipleElements(t *testing.T) {
	whitelist = []string{"3456789876543234567", "45678998765465"}
	whiteListed := isWhitelisted(45678998765465)
	assert.Assert(t, whiteListed)
}

func TestIsWhitelistedFalseMultipleElements(t *testing.T) {
	whitelist = []string{"3456789876543234567", "45678998765465"}
	whiteListed := isWhitelisted(87654334567)
	assert.Assert(t, !whiteListed)
}