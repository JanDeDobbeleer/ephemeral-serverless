package main

import (
	"errors"
	"github.com/ChimeraCoder/anaconda"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"net/url"
	"os"
	"testing"
)

// Mocked AnacondaClient to test
type MockedTwitterClient struct {
	mock.Mock
}

func (m *MockedTwitterClient) GetSelf(v url.Values) (u anaconda.User, err error) {
	args := m.Called(v)
	return args.Get(0).(anaconda.User), args.Error(1)
}

func (m *MockedTwitterClient) GetUserTimeline(v url.Values) (timeline []anaconda.Tweet, err error) {
	args := m.Called(nil)
	return args.Get(0).([]anaconda.Tweet), args.Error(1)
}

func (m *MockedTwitterClient) GetSearch(queryString string, v url.Values) (sr anaconda.SearchResponse, err error) {
	args := m.Called(queryString, v)
	return args.Get(0).(anaconda.SearchResponse), args.Error(1)
}

func (m *MockedTwitterClient) DeleteTweet(id int64, trimUser bool) (tweet anaconda.Tweet, err error) {
	args := m.Called(id, trimUser)
	return args.Get(0).(anaconda.Tweet), args.Error(1)
}

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

func TestGetTimelineRandomTweets(t *testing.T) {
	client := new(MockedTwitterClient)
	numberOftweets := rand.Intn(200)
	var tweets []anaconda.Tweet
	for i := 0; i < numberOftweets; i++ {
		tweets = append(tweets, anaconda.Tweet{
			Id: rand.Int63(),
		})
	}
	client.On("GetUserTimeline", nil).Return(tweets, nil)
	timeline, err := getTimeline(client)
	assert.EqualValues(t, tweets, timeline)
	assert.Nil(t, err)
}

func TestGetTimelineError(t *testing.T) {
	client := new(MockedTwitterClient)
	tweets := make([]anaconda.Tweet, 0)
	expectedError := errors.New("emit macho dwarf: elf header corrupted")
	client.On("GetUserTimeline", nil).Return(tweets, expectedError)
	timeline, err := getTimeline(client)
	assert.EqualValues(t, tweets, timeline)
	assert.Equal(t, err, expectedError)
}

func TestGetTimelineNoTweets(t *testing.T) {
	client := new(MockedTwitterClient)
	tweets := make([]anaconda.Tweet, 0)
	client.On("GetUserTimeline", nil).Return(tweets, nil)
	timeline, err := getTimeline(client)
	assert.EqualValues(t, tweets, timeline)
	assert.Nil(t, err)
}
