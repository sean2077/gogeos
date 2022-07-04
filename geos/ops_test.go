package geos

import (
	"fmt"
	"testing"
)

func TestSplitPolygonWithLine(t *testing.T) {
	type args struct {
		polygon string
		line    string
	}
	tests := []struct {
		input   args
		want    []string
		wantErr bool
	}{
		{
			input: args{
				polygon: "POLYGON ((20 20, 20 10, 10 10, 10 20, 20 20))",
				line:    "LINESTRING (20 20, 10 10)",
			},
			want: []string{
				"POLYGON ((20 20, 20 10, 10 10, 20 20))",
				"POLYGON ((10 10, 10 20, 20 20, 10 10))",
			},
			wantErr: false,
		},
		{
			input: args{
				polygon: "POLYGON ((20 20, 20 10, 10 10, 10 20, 20 20))",
				line:    "LINESTRING (0 0, 1 1)",
			},
			want: []string{
				"POLYGON ((20 20, 20 10, 10 10, 10 20, 20 20))",
			},
			wantErr: false,
		},
		{
			input: args{
				polygon: "POLYGON ((20 20, 20 10, 10 10, 10 20, 20 20))",
				line:    "LINESTRING (5 10, 25 20)",
			},
			want: []string{
				"POLYGON ((20 20, 20 17.5, 10 12.5, 10 20, 20 20))",
				"POLYGON ((20 17.5, 20 10, 10 10, 10 12.5, 20 17.5))",
			},
			wantErr: false,
		},
	}
	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("case%d", i), func(t *testing.T) {
				poly, _ := FromWKT(tt.input.polygon)
				line, _ := FromWKT(tt.input.line)
				var want []*Geometry
				for _, wkt := range tt.want {
					g, _ := FromWKT(wkt)
					want = append(want, g)
				}
				got, err := SplitPolygonWithLine(poly, line)
				if (err != nil) != tt.wantErr {
					t.Errorf("SplitPolygonWithLine() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if equal, err := compare(got, want); err != nil || !equal {
					t.Errorf("SplitPolygonWithLine() got = %v, want %v", got, want)
				}
			},
		)
	}
}

func compare(mp interface{}, polys []*Geometry) (bool, error) {
	switch t := mp.(type) {
	case *Geometry:
		npolys, err := t.NGeometry()
		if err != nil {
			return false, err
		}
		if npolys != len(polys) {
			return false, nil
		}
		for i := range polys {
			p, err := t.Geometry(i)
			if err != nil {
				return false, err
			}
			equal, err := p.Equals(polys[i])
			if err != nil {
				return false, err
			}
			if !equal {
				return false, nil
			}
		}
		return true, nil

	case []*Geometry:
		if len(t) != len(polys) {
			return false, nil
		}
		for i := range polys {
			equal, err := t[i].Equals(polys[i])
			if err != nil {
				return false, err
			}
			if !equal {
				return false, nil
			}
		}
		return true, nil
	}
	return false, nil
}
