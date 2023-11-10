package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckArgs_ValidArguments(t *testing.T) {
	host, port := "localhost", "8080"
	result := checkArgs(host, port)
	assert.True(t, result, "Expected checkArgs to return true for valid arguments")
}

func TestCheckArgs_ValidHostDomain(t *testing.T) {
	host, port := "google.com", "8080"
	result := checkArgs(host, port)
	assert.True(t, result, "Expected checkArgs to return false for invalid host")
}

func TestCheckArgs_ValidHostIP(t *testing.T) {
	host, port := "8.8.8.8", "8080"
	result := checkArgs(host, port)
	assert.True(t, result, "Expected checkArgs to return false for invalid host")
}

func TestCheckArgs_InvalidArgumentCount(t *testing.T) {
	host, port := "127.0.0.1", ""
	result := checkArgs(host, port)
	assert.False(t, result, "Expected checkArgs to return false for invalid argument count")
}

func TestCheckArgs_InvalidPort(t *testing.T) {
	host, port := "exemple.com", "abs"
	result := checkArgs(host, port)
	assert.False(t, result, "Expected checkArgs to return false for invalid port")
}

func TestCheckArgs_InvalidPort2(t *testing.T) {
	host, port := "exemple.com", "80808080"
	result := checkArgs(host, port)
	assert.False(t, result, "Expected checkArgs to return false for invalid port")
}

func TestCheckArgs_InvalidHost(t *testing.T) {
	host, port := "unknownhost", "8080"
	result := checkArgs(host, port)
	assert.False(t, result, "Expected checkArgs to return false for invalid host")
}

func TestCheckArgs_InvalidHostDomain(t *testing.T) {
	host, port := "somehostinvalid.ru", "8080"
	result := checkArgs(host, port)
	assert.False(t, result, "Expected checkArgs to return false for invalid host")
}

func TestCheckArgs_InvalidHostIP(t *testing.T) {
	host, port := "256.0.0.1", "8080"
	result := checkArgs(host, port)
	assert.False(t, result, "Expected checkArgs to return false for invalid host")
}
