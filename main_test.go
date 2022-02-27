package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestFetchIAM(t *testing.T) {
    // connect to the AWS API with the access key pair
    assert.Equal(t, nil, FetchIAM(), "Failed to fetch IAM profiles")
    // find out IAM access keys over N hours
    // save a file which contains the suggested result
}