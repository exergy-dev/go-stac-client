// CQL2 filter helpers for STAC API search.
//
// This file provides a thin convenience layer over the
// github.com/exergy-dev/go-cql2 library so STAC users can build, encode, and
// inject CQL2 filters into SearchParams without boilerplate.
//
// Example:
//
//	import (
//	    cql2 "github.com/exergy-dev/go-cql2"
//	    stacclient "github.com/robert-malhotra/go-stac-client/pkg/client"
//	)
//
//	expr := cql2.And(
//	    cql2.Lt("eo:cloud_cover", 10),
//	    cql2.SIntersects("geometry", cql2.Geom(/* a Point or Polygon */)),
//	)
//	params := stacclient.SearchParams{
//	    Collections: []string{"sentinel-2-l2a"},
//	    Filter:      stacclient.MustCQL2JSON(expr),
//	    FilterLang:  stacclient.FilterLangCQL2JSON,
//	}
//	for item, err := range cli.Search(ctx, params) {
//	    // ...
//	}
package client

import (
	"encoding/json"
	"fmt"

	cql2 "github.com/exergy-dev/go-cql2"
	cql2json "github.com/exergy-dev/go-cql2/json"
	cql2text "github.com/exergy-dev/go-cql2/text"
)

// Filter language identifiers for SearchParams.FilterLang.
const (
	FilterLangCQL2JSON = "cql2-json"
	FilterLangCQL2Text = "cql2-text"
)

// CQL2JSON encodes a CQL2 expression as a JSON-encoded filter suitable for
// SearchParams.Filter when FilterLang is "cql2-json".
//
// The returned json.RawMessage is round-trip-stable: it embeds verbatim into
// the search-request body without re-encoding.
func CQL2JSON(expr cql2.Expr) (json.RawMessage, error) {
	if expr.Node() == nil {
		return nil, fmt.Errorf("stac: empty CQL2 expression")
	}
	b, err := cql2json.Encode(expr.Node())
	if err != nil {
		return nil, fmt.Errorf("stac: encode cql2-json: %w", err)
	}
	return json.RawMessage(b), nil
}

// MustCQL2JSON is the panic-on-error variant of CQL2JSON.
func MustCQL2JSON(expr cql2.Expr) json.RawMessage {
	out, err := CQL2JSON(expr)
	if err != nil {
		panic(err)
	}
	return out
}

// CQL2Text encodes a CQL2 expression as the CQL2-text string suitable for
// SearchParams.Filter when FilterLang is "cql2-text".
func CQL2Text(expr cql2.Expr) (string, error) {
	if expr.Node() == nil {
		return "", fmt.Errorf("stac: empty CQL2 expression")
	}
	b, err := cql2text.Encode(expr.Node())
	if err != nil {
		return "", fmt.Errorf("stac: encode cql2-text: %w", err)
	}
	return string(b), nil
}

// MustCQL2Text is the panic-on-error variant of CQL2Text.
func MustCQL2Text(expr cql2.Expr) string {
	out, err := CQL2Text(expr)
	if err != nil {
		panic(err)
	}
	return out
}
