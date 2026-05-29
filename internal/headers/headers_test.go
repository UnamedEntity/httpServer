package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Unit test for headers
func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra white space
	headers = NewHeaders()
	data = []byte("Host:         localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 31, n)
	assert.False(t, done)

	// Test: Valid 2 headers with exsiting headers
	headers = NewHeaders()
	data = []byte("Host:         localhost:42069\r\nLand: blalbla\r\n\r\n")
	n, done, err = headers.Parse(data)
	n, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "blalbla", headers["land"])
	assert.Equal(t, 15, n)
	assert.False(t, done)

	// Test: Valid Done
	headers = NewHeaders()
	data = []byte("Host:         localhost:42069\r\nLand: blalbla\r\n\r\n\r\n")
	n, done, err = headers.Parse(data)
	total := n
	n, done, err = headers.Parse(data[n:])
	total += n
	n, done, err = headers.Parse(data[total:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 0, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid chracters
	headers = NewHeaders()
	data = []byte("Host<>><><><><><><><: localhost<><><:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

}
