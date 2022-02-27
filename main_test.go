package main

import (
	"time"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestFetchIAM(t *testing.T) {
    // connect to the AWS API with the access key pair
	fetched, err := FetchIAM()
    assert.Equal(t, nil, err, "Failed to fetch IAM profiles")
    // find out IAM access keys over N hours
	hours := time.Duration(time.Hour * 720)
	assert.Equal(t, 1, len(CheckProfileExpired(hours, fetched)), "Counts of targeted IAM profiles are not the same as expected")
    // save a file which contains the suggested result
}