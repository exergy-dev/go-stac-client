package client

import (
	"encoding/json"
	"testing"

	cql2 "github.com/exergy-dev/go-cql2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCQL2JSON(t *testing.T) {
	expr := cql2.And(
		cql2.Lt("eo:cloud_cover", 10),
		cql2.Eq("collection", "sentinel-2-l2a"),
	)
	out, err := CQL2JSON(expr)
	require.NoError(t, err)

	// Round-trip the JSON back through go-cql2.
	var generic any
	require.NoError(t, json.Unmarshal(out, &generic))

	root := generic.(map[string]any)
	assert.Equal(t, "and", root["op"])
	args, ok := root["args"].([]any)
	require.True(t, ok)
	assert.Len(t, args, 2)
}

func TestCQL2Text(t *testing.T) {
	expr := cql2.Lt("eo:cloud_cover", 10)
	out, err := CQL2Text(expr)
	require.NoError(t, err)
	assert.Contains(t, out, "eo:cloud_cover")
	assert.Contains(t, out, "10")
}

func TestCQL2JSONEmptyExprErrors(t *testing.T) {
	var empty cql2.Expr
	_, err := CQL2JSON(empty)
	assert.Error(t, err)
}

func TestCQL2JSONIntegratesWithSearchParams(t *testing.T) {
	expr := cql2.Eq("id", "abc")
	filter := MustCQL2JSON(expr)
	params := SearchParams{
		Filter:     filter,
		FilterLang: FilterLangCQL2JSON,
	}
	body, err := json.Marshal(params)
	require.NoError(t, err)

	var decoded map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(body, &decoded))
	require.Contains(t, decoded, "filter")
	// The filter must round-trip as a JSON object, not a string.
	assert.Equal(t, byte('{'), decoded["filter"][0])
	assert.Equal(t, "\"cql2-json\"", string(decoded["filter-lang"]))
}
