go-geojson
==========

go-geojson provides a fluent [Go](https://go.dev/) API for building and marshalling [GeoJSON](https://geojson.org/) objects.

This library was written as an alternative to [github.com/paulmach/go.geojson](https://github.com/paulmach/go.geojson),
which isn't as nice to use when creating objects (e.g., one can't fluently set properties when creating a feature,
it doesn't guide the user as to how to create positions or how many positions need to be in a line string or linear ring, etc.).
github.com/paulmach/go.geojson currently supports more features though, in particular unmarshaling and bounding boxes.

To get the library:
```bash
go get github.com/barnardb/go-geojson
```
