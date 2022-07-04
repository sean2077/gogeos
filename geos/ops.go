package geos

import (
	"log"
)

// SplitPolygonWithLine returns parts of the polygon split by line string.
func SplitPolygonWithLine(poly, splitter *Geometry) ([]*Geometry, error) {
	if t, err := poly.Type(); err != nil || t != POLYGON {
		log.Fatal("First argument must be a Polygon")
	}
	if t, err := splitter.Type(); err != nil || t != LINESTRING {
		log.Fatal("Second argument must be a LineString")
	}

	boundary, err := poly.Boundary()
	if err != nil {
		return nil, err
	}
	union, err := boundary.Union(splitter)
	if err != nil {
		return nil, err
	}
	// greatly improves split performance for big geometries with many
	// holes (the following contains checks) with minimal overhead
	// for common cases
	ppoly := poly.Prepare()

	nlines, _ := union.NGeometry()
	lines := make([]*Geometry, nlines)
	for i := 0; i < nlines; i++ {
		lines[i], _ = union.Geometry(i)
	}
	polys, err := Polygonize(lines)
	if err != nil {
		return nil, err
	}

	var res []*Geometry

	npolys, _ := polys.NGeometry()
	for i := 0; i < npolys; i++ {
		pg, _ := polys.Geometry(i)
		rpt, err := pg.PointOnSurface()
		if err != nil {
			continue
		}
		contain, err := ppoly.Contains(rpt)
		if err == nil && contain {
			res = append(res, pg)
		}
	}

	return res, nil
}
