package main

import (
	"errors"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"net/url"
	"os"
	"testing"
	"time"
)

// Mocked AnacondaClient to test
type MockedTwitterClient struct {
	mock.Mock
}

func (m *MockedTwitterClient) GetSelf(v url.Values) (u anaconda.User, err error) {
	args := m.Called(nil)
	return args.Get(0).(anaconda.User), args.Error(1)
}

func (m *MockedTwitterClient) GetUserTimeline(v url.Values) (timeline []anaconda.Tweet, err error) {
	args := m.Called(nil)
	return args.Get(0).([]anaconda.Tweet), args.Error(1)
}

func (m *MockedTwitterClient) GetSearch(queryString string, v url.Values) (sr anaconda.SearchResponse, err error) {
	args := m.Called(queryString, nil)
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

func TestGetRepliesForTweetRandom(t *testing.T) {
	shouldMatch := rand.Intn(200)
	shouldNotMatch := rand.Intn(200)
	tweetID := rand.Int63()
	client := new(MockedTwitterClient)
	me := anaconda.User{
		ScreenName: "elonmusk",
	}
	queryString := fmt.Sprintf("to:%s", me.ScreenName)
	client.On("GetSelf", nil).Return(me, nil)
	var tweets []anaconda.Tweet
	for i := 0; i < shouldNotMatch; i++ {
		tweets = append(tweets, anaconda.Tweet{
			Id: rand.Int63(),
		})
	}
	for i := 0; i < shouldMatch; i++ {
		tweets = append(tweets, anaconda.Tweet{
			Id:                rand.Int63(),
			InReplyToStatusID: tweetID,
		})
	}
	searchResponse := anaconda.SearchResponse{
		Statuses: tweets,
	}
	client.On("GetSearch", queryString, nil).Return(searchResponse, nil)
	replies := getRepliesForTweet(client, tweetID)
	assert.Len(t, replies, shouldMatch)
}

func TestGetRepliesForTweetSelfError(t *testing.T) {
	tweetID := rand.Int63()
	client := new(MockedTwitterClient)
	client.On("GetSelf", nil).Return(anaconda.User{}, errors.New("i'm a little teapot short and stout"))
	replies := getRepliesForTweet(client, tweetID)
	assert.Len(t, replies, 0)
}

func TestGetRepliesForTweetSearchError(t *testing.T) {
	tweetID := rand.Int63()
	client := new(MockedTwitterClient)
	me := anaconda.User{
		ScreenName: "elonmusk",
	}
	queryString := fmt.Sprintf("to:%s", me.ScreenName)
	client.On("GetSelf", nil).Return(me, nil)
	err := errors.New("i'm a little teapot short and stout")
	client.On("GetSearch", queryString, nil).Return(anaconda.SearchResponse{}, err)
	replies := getRepliesForTweet(client, tweetID)
	assert.Len(t, replies, 0)
}

func getHasOngoingInteractions(hasActiveInteractions bool, useValidDates bool) bool {
	shouldMatch := rand.Intn(200)
	shouldNotMatch := rand.Intn(200)
	tweetID := rand.Int63()
	interactionAgeLimitHours := rand.Intn(200)
	client := new(MockedTwitterClient)
	me := anaconda.User{
		ScreenName: "elonmusk",
	}
	queryString := fmt.Sprintf("to:%s", me.ScreenName)
	client.On("GetSelf", nil).Return(me, nil)
	var tweets []anaconda.Tweet
	for i := 0; i < shouldNotMatch; i++ {
		tweets = append(tweets, anaconda.Tweet{
			Id: rand.Int63(),
		})
	}
	var replies []anaconda.Tweet
	for i := 0; i < shouldMatch; i++ {
		createdAt := ""
		if useValidDates {
			createdAt = time.Now().Add(-201 * time.Hour).Format(time.RubyDate)
		}
		replies = append(replies, anaconda.Tweet{
			Id:                rand.Int63(),
			InReplyToStatusID: tweetID,
			CreatedAt:         createdAt,
		})
	}
	if hasActiveInteractions {
		shouldBeActive := rand.Intn(shouldMatch)
		for i := 0; i < shouldBeActive; i++ {
			randomHour := rand.Intn(interactionAgeLimitHours)
			createdAt := time.Now().Add(time.Duration(-randomHour) * time.Hour)
			replies[i].CreatedAt = createdAt.Format(time.RubyDate)
		}
	}
	tweets = append(tweets, replies...)
	searchResponse := anaconda.SearchResponse{
		Statuses: tweets,
	}
	client.On("GetSearch", queryString, nil).Return(searchResponse, nil)
	interactionAgeLimitFormat := fmt.Sprintf("%dh", interactionAgeLimitHours)
	interactionAgeLimit, _ := time.ParseDuration(interactionAgeLimitFormat)
	return hasOngoingInteractions(client, tweetID, interactionAgeLimit)
}

func TestHasOngoingInteractionsTrue(t *testing.T) {
	hasOngoingInteractions := getHasOngoingInteractions(true, true)
	assert.True(t, hasOngoingInteractions)
}

func TestHasOngoingInteractionsFalse(t *testing.T) {
	hasOngoingInteractions := getHasOngoingInteractions(false, true)
	assert.False(t, hasOngoingInteractions)
}

func TestHasOngoingInteractionsInvalidDatesTrue(t *testing.T) {
	hasOngoingInteractions := getHasOngoingInteractions(true, false)
	assert.True(t, hasOngoingInteractions)
}

func TestHasOngoingInteractionsInvalidDatesFalse(t *testing.T) {
	hasOngoingInteractions := getHasOngoingInteractions(false, false)
	assert.False(t, hasOngoingInteractions)
}
