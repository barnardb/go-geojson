// Package geojson provides utilities for programatically creating GeoJSON output.
// GeoJSON is standardized in [IEEE RFC 7946: The GeoJSON Format](https://www.rfc-editor.org/rfc/rfc7946).
package geojson

import (
	"encoding/json"
	"fmt"
)

// Object represents a [GeoJSON object](https://www.rfc-editor.org/rfc/rfc7946#section-3).
type Object interface {

	// Type returns the type of this GeoJSON object.
	//
	// The types are introduced in [§1.4 Definitions](https://www.rfc-editor.org/rfc/rfc7946#section-1.4)
	// as being the geometry types "Point", "MultiPoint", "LineString", "MultiLineString", "Polygon", "MultiPolygon", and "GeometryCollection",
	// as well as "Feature" and "FeatureCollection".
	// The type required by each object is documented in the subsections of [§2 GeoJSON Object](https://www.rfc-editor.org/rfc/rfc7946#section-2) specific to each object.
	// They are again enumerated in [§7 GeoJSON Types Are Not Extensible](https://www.rfc-editor.org/rfc/rfc7946#section-7),
	// which emphasises that other types are not allowed.
	Type() string

	// // BBox() *BoundingBox
}

var _ Object = (Geometry)(nil)
var _ Object = (*Feature)(nil)
var _ Object = (*FeatureCollection)(nil)

// ToText converts an Object to [GeoJSON text](https://www.rfc-editor.org/rfc/rfc7946#section-2).
func ToText(o Object) ([]byte, error) {
	return json.Marshal(o)
}

// Position represents a [GeoJSON position](https://www.rfc-editor.org/rfc/rfc7946#section-3.1.1).
type Position struct {
	Longitude, Latitude float64
	// Altitude           *float64
}

// MarshalJSON implements json.Marshaler.
func (p Position) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%f,%f]", p.Longitude, p.Latitude)), nil
}

// Geometry represents a [GeoJSON Geometry](https://www.rfc-editor.org/rfc/rfc7946#section-3.1).
type Geometry interface {
	Object
	ToFeature(properties map[string]any) Feature
	// Coordinates() []any
}

var _ Geometry = (*Point)(nil)
var _ Geometry = (*MultiPoint)(nil)
var _ Geometry = (*LineString)(nil)
var _ Geometry = (*MultiLineString)(nil)
var _ Geometry = (*Polygon)(nil)
var _ Geometry = (*MultiPolygon)(nil)
var _ Geometry = (*GeometryCollection)(nil)

type geometryWithCoordinates[C any] struct {
	Coordinates C `json:"coordinates"`
}

// LineStringCoordinates represents the coordinates of a LineString, and one of the elements in the coordinates of a MultiLineString.
type LineStringCoordinates []Position

// NewLineStringCoordinates creates a new LineStringCoordinates with the given positions.
func NewLineStringCoordinates(c0, c1 Position, cn ...Position) LineStringCoordinates {
	coordinates := make(LineStringCoordinates, 0, 2+len(cn))
	return append(append(coordinates, c0, c1), cn...)
}

// LinearRing represents the coordinates of a single linear ring element of a Polygon coordinate array.
type LinearRing []Position

// NewLinearRing creates a new LinearRing with the given positions.
func NewLinearRing(c0, c1, c2, c3 Position, cn ...Position) LinearRing {
	coordinates := make(LinearRing, 0, 4+len(cn))
	coordinates = append(append(coordinates, c0, c1, c2, c3), cn...)
	if coordinates[len(coordinates)-1] != c0 {
		panic(fmt.Sprintf("start position %#v doesn't match end position %#v", c0, coordinates[len(coordinates)-1]))
	}
	return coordinates
}

// PolygonCoordinates represents the coordinates of a Polygon, and one of the elements in the coordinates of a MultiPolygon.
type PolygonCoordinates []LinearRing

// NewPolygonCoordinates creates a new PolygonCoordinates with the given linear rings.
func NewPolygonCoordinates(outerBoundary LinearRing, holes ...LinearRing) PolygonCoordinates {
	coordinates := make(PolygonCoordinates, 0, 1+len(holes))
	return append(append(coordinates, outerBoundary), holes...)
}

// Point represents a [GeoJSON Point](https://www.rfc-editor.org/rfc/rfc7946#section-3.1.2).
type Point geometryWithCoordinates[Position]

// NewPoint creates a new Point with the given position.
func NewPoint(p Position) *Point {
	return &Point{p}
}

func marshalGeoJSONGeometryWithCoordinates[C any](g Geometry, c C) ([]byte, error) {
	return json.Marshal(struct {
		Type        string `json:"type"`
		Coordinates C      `json:"coordinates"`
	}{
		Type:        g.Type(),
		Coordinates: c,
	})
}

// Type returns "Point", implementing interface Object.
func (g *Point) Type() string {
	return "Point"
}

// MarhsalJSON implements json.Marshaler.
func (g *Point) ToFeature(properties map[string]any) Feature { return Feature{g, properties} }

// MarshalJSON implements json.Marshaler.
func (g *Point) MarshalJSON() ([]byte, error) {
	return marshalGeoJSONGeometryWithCoordinates(g, g.Coordinates)
}

// MultiPoint represents a [GeoJSON MultiPoint](https://www.rfc-editor.org/rfc/rfc7946#section-3.1.3).
type MultiPoint geometryWithCoordinates[[]Position]

// NewMultiPoint creates a new MultiPoint with the given position.
func NewMultiPoint(pn ...Position) *MultiPoint {
	return &MultiPoint{pn}
}

// Type returns "MultiPoint", implementing interface Object.
func (g *MultiPoint) Type() string {
	return "MultiPoint"
}
func (g *MultiPoint) ToFeature(properties map[string]any) Feature { return Feature{g, properties} }

// MarshalJSON implements json.Marshaler.
func (g *MultiPoint) MarshalJSON() ([]byte, error) {
	return marshalGeoJSONGeometryWithCoordinates(g, g.Coordinates)
}

// LineString represents a [GeoJSON LineString](https://www.rfc-editor.org/rfc/rfc7946#section-3.1.4).
type LineString geometryWithCoordinates[LineStringCoordinates]

// NewLineString creates a new LineString with the given positions.
func NewLineString(c0, c1 Position, cn ...Position) *LineString {
	return &LineString{NewLineStringCoordinates(c0, c1, cn...)}
}

// Type returns "LineString", implementing interface Object.
func (g *LineString) Type() string {
	return "LineString"
}
func (g *LineString) ToFeature(properties map[string]any) Feature { return Feature{g, properties} }

// MarshalJSON implements json.Marshaler.
func (g *LineString) MarshalJSON() ([]byte, error) {
	return marshalGeoJSONGeometryWithCoordinates(g, g.Coordinates)
}

// MultiLineString represents a [GeoJSON MultiLineString](https://www.rfc-editor.org/rfc/rfc7946#section-3.1.5).
type MultiLineString geometryWithCoordinates[[]LineStringCoordinates]

// NewMultiLineString creates a new MultiLineString with the given LineStringCoordinates.
func NewMultiLineString(coordinates ...LineStringCoordinates) *MultiLineString {
	return &MultiLineString{coordinates}
}

// Type returns "MultiLineString", implementing interface Object.
func (g *MultiLineString) Type() string {
	return "MultiLineString"
}
func (g *MultiLineString) ToFeature(properties map[string]any) Feature { return Feature{g, properties} }

// MarshalJSON implements json.Marshaler.
func (g *MultiLineString) MarshalJSON() ([]byte, error) {
	return marshalGeoJSONGeometryWithCoordinates(g, g.Coordinates)
}

// Polygon represents a [GeoJSON Polygon](https://www.rfc-editor.org/rfc/rfc7946#section-3.1.6).
type Polygon geometryWithCoordinates[PolygonCoordinates]

// NewPolygon creates a new Polygon with the given linear rings.
func NewPolygon(outerBoundary LinearRing, holes ...LinearRing) *Polygon {
	return &Polygon{NewPolygonCoordinates(outerBoundary, holes...)}
}

// Type returns "Polygon", implementing interface Object.
func (g *Polygon) Type() string {
	return "Polygon"
}
func (g *Polygon) ToFeature(properties map[string]any) Feature { return Feature{g, properties} }

// MarshalJSON implements json.Marshaler.
func (g *Polygon) MarshalJSON() ([]byte, error) {
	return marshalGeoJSONGeometryWithCoordinates(g, g.Coordinates)
}

// MultiPolygon represents a [GeoJSON MultiPolygon](https://www.rfc-editor.org/rfc/rfc7946#section-3.1.7).
type MultiPolygon geometryWithCoordinates[[]PolygonCoordinates]

// NewMultiPolygon creates a new MultiPolygon from the given PolygoneCoordinates.
func NewMultiPolygon(coordinates ...PolygonCoordinates) *MultiPolygon {
	return &MultiPolygon{coordinates}
}

// Type returns "MultiPolygon", implementing interface Object.
func (g *MultiPolygon) Type() string {
	return "MultiPolygon"
}
func (g *MultiPolygon) ToFeature(properties map[string]any) Feature { return Feature{g, properties} }

// MarshalJSON implements json.Marshaler.
func (g *MultiPolygon) MarshalJSON() ([]byte, error) {
	return marshalGeoJSONGeometryWithCoordinates(g, g.Coordinates)
}

// GeometryCollection represents a [GeoJSON GeometryCollection](https://www.rfc-editor.org/rfc/rfc7946#section-3.1.8).
type GeometryCollection []Geometry

// NewGeometryCollection creates a new GeometryCollection from the given Geomtries.
func NewGeometryCollection(geometries ...Geometry) *GeometryCollection {
	return (*GeometryCollection)(&geometries)
}

// Type returns "GeometryCollection", implementing interface Object.
func (g *GeometryCollection) Type() string {
	return "GeometryCollection"
}
func (g *GeometryCollection) ToFeature(properties map[string]any) Feature {
	return Feature{g, properties}
}

// Feature represents a [GeoJSON Feature](https://www.rfc-editor.org/rfc/rfc7946#section-3.2).
type Feature struct {
	Geometry   Geometry       `json:"geometry"`
	Properties map[string]any `json:"properties"`
	// ID string | float64 `json:"id"`
}

// NewFeature creates a new Feature from the given Geometry and properties.
func NewFeature(geometry Geometry, properties map[string]any) Feature {
	return Feature{geometry, properties}
}

// Type returns "Feature", implementing interface Object.
func (f Feature) Type() string {
	return "Feature"
}

// MarshalJSON implements json.Marshaler.
func (f *Feature) MarshalJSON() ([]byte, error) {
	type Raw Feature
	return json.Marshal(struct {
		Type string `json:"type"`
		*Raw
	}{
		Type: f.Type(),
		Raw:  (*Raw)(f),
	})
}

// FeatureCollection represents a [GeoJSON FeatureCollection](https://www.rfc-editor.org/rfc/rfc7946#section-3.3).
type FeatureCollection []Feature

// NewFeatureCollection creates a new FeatureCollection from the given Features.
func NewFeatureCollection(features ...Feature) FeatureCollection {
	return features
}

// Type returns "FeatureCollection", implementing interface Object.
func (c *FeatureCollection) Type() string {
	return "FeatureCollection"
}

// Type returns "FeatureCollection", implementing interface Object.
func (c FeatureCollection) With(features ...Feature) FeatureCollection {
	return append(c, features...)
}

// MarshalJSON implements json.Marshaler.
func (c *FeatureCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string    `json:"type"`
		Features []Feature `json:"features"`
	}{
		Type:     c.Type(),
		Features: *c,
	})
}
