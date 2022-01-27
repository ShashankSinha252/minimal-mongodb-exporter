package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDSN(t *testing.T) {
	host := "localhost"
	port := 27017

	dsn := createDSN(host, "", "", port)
	expect := fmt.Sprintf("mongodb://%s:%d", host, port)
	assert.Equal(t, expect, dsn)
}
