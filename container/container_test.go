package container

import (
	"image"
	"log"
	"testing"

	"github.com/mum4k/termdash/canvas"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/draw"
	"github.com/mum4k/termdash/terminal/faketerm"
)

// Example demonstrates how to use the Container API.
func Example() {
	New(
		/* terminal = */ nil,
		SplitVertical(
			Left(
				SplitHorizontal(
					Top(
						Border(draw.LineStyleLight),
					),
					Bottom(
						SplitHorizontal(
							Top(
								Border(draw.LineStyleLight),
							),
							Bottom(
								Border(draw.LineStyleLight),
							),
						),
					),
				),
			),
			Right(
				Border(draw.LineStyleLight),
			),
		),
	)

	// TODO(mum4k): Allow splits on different ratios.
	// TODO(mum4k): Include an example with a widget.
}

// mustCanvas returns a new canvas or panics.
func mustCanvas(area image.Rectangle) *canvas.Canvas {
	cvs, err := canvas.New(area)
	if err != nil {
		log.Fatalf("canvas.New => unexpected error: %v", err)
	}
	return cvs
}

// mustBox draws box on the canvas or panics.
func mustBox(c *canvas.Canvas, box image.Rectangle, ls draw.LineStyle, opts ...cell.Option) {
	if err := draw.Box(c, box, ls, opts...); err != nil {
		log.Fatalf("draw.Box => unexpected error: %v", err)
	}
}

// mustApply applies the canvas on the terminal or panics.
func mustApply(c *canvas.Canvas, t *faketerm.Terminal) {
	if err := c.Apply(t); err != nil {
		log.Fatalf("canvas.Apply => unexpected error: %v", err)
	}
}

func TestDraw(t *testing.T) {
	tests := []struct {
		desc      string
		termSize  image.Point
		container func(ft *faketerm.Terminal) *Container
		want      func(size image.Point) *faketerm.Terminal
		wantErr   bool
	}{
		{
			desc:     "empty container",
			termSize: image.Point{10, 10},
			container: func(ft *faketerm.Terminal) *Container {
				return New(ft)
			},
			want: func(size image.Point) *faketerm.Terminal {
				return faketerm.MustNew(size)
			},
		},
		{
			desc:     "container with a border",
			termSize: image.Point{10, 10},
			container: func(ft *faketerm.Terminal) *Container {
				return New(
					ft,
					Border(draw.LineStyleLight),
				)
			},
			want: func(size image.Point) *faketerm.Terminal {
				ft := faketerm.MustNew(size)
				cvs := mustCanvas(image.Rect(0, 0, 10, 10))
				mustBox(cvs, image.Rect(0, 0, 10, 10), draw.LineStyleLight)
				mustApply(cvs, ft)
				return ft
			},
		},
		{
			desc:     "horizontal split, children have borders",
			termSize: image.Point{10, 10},
			container: func(ft *faketerm.Terminal) *Container {
				return New(
					ft,
					SplitHorizontal(
						Top(
							Border(draw.LineStyleLight),
						),
						Bottom(
							Border(draw.LineStyleLight),
						),
					),
				)
			},
			want: func(size image.Point) *faketerm.Terminal {
				ft := faketerm.MustNew(size)
				cvs := mustCanvas(image.Rect(0, 0, 10, 10))
				mustBox(cvs, image.Rect(0, 0, 10, 5), draw.LineStyleLight)
				mustBox(cvs, image.Rect(0, 5, 10, 10), draw.LineStyleLight)
				mustApply(cvs, ft)
				return ft
			},
		},
		{
			desc:     "horizontal split, parent and children have borders",
			termSize: image.Point{10, 10},
			container: func(ft *faketerm.Terminal) *Container {
				return New(
					ft,
					Border(draw.LineStyleLight),
					SplitHorizontal(
						Top(
							Border(draw.LineStyleLight),
						),
						Bottom(
							Border(draw.LineStyleLight),
						),
					),
				)
			},
			want: func(size image.Point) *faketerm.Terminal {
				ft := faketerm.MustNew(size)
				cvs := mustCanvas(image.Rect(0, 0, 10, 10))
				mustBox(cvs, image.Rect(0, 0, 10, 10), draw.LineStyleLight)
				mustBox(cvs, image.Rect(1, 1, 9, 5), draw.LineStyleLight)
				mustBox(cvs, image.Rect(1, 5, 9, 9), draw.LineStyleLight)
				mustApply(cvs, ft)
				return ft
			},
		},
		{
			desc:     "vertical split, children have borders",
			termSize: image.Point{10, 10},
			container: func(ft *faketerm.Terminal) *Container {
				return New(
					ft,
					SplitVertical(
						Left(
							Border(draw.LineStyleLight),
						),
						Right(
							Border(draw.LineStyleLight),
						),
					),
				)
			},
			want: func(size image.Point) *faketerm.Terminal {
				ft := faketerm.MustNew(size)
				cvs := mustCanvas(image.Rect(0, 0, 10, 10))
				mustBox(cvs, image.Rect(0, 0, 5, 10), draw.LineStyleLight)
				mustBox(cvs, image.Rect(5, 0, 10, 10), draw.LineStyleLight)
				mustApply(cvs, ft)
				return ft
			},
		},
		{
			desc:     "vertical split, parent and children have borders",
			termSize: image.Point{10, 10},
			container: func(ft *faketerm.Terminal) *Container {
				return New(
					ft,
					Border(draw.LineStyleLight),
					SplitVertical(
						Left(
							Border(draw.LineStyleLight),
						),
						Right(
							Border(draw.LineStyleLight),
						),
					),
				)
			},
			want: func(size image.Point) *faketerm.Terminal {
				ft := faketerm.MustNew(size)
				cvs := mustCanvas(image.Rect(0, 0, 10, 10))
				mustBox(cvs, image.Rect(0, 0, 10, 10), draw.LineStyleLight)
				mustBox(cvs, image.Rect(1, 1, 5, 9), draw.LineStyleLight)
				mustBox(cvs, image.Rect(5, 1, 9, 9), draw.LineStyleLight)
				mustApply(cvs, ft)
				return ft
			},
		},
		{
			desc:     "multi level split",
			termSize: image.Point{10, 11},
			container: func(ft *faketerm.Terminal) *Container {
				return New(
					ft,
					SplitVertical(
						Left(
							SplitHorizontal(
								Top(
									Border(draw.LineStyleLight),
								),
								Bottom(
									SplitHorizontal(
										Top(
											Border(draw.LineStyleLight),
										),
										Bottom(
											Border(draw.LineStyleLight),
										),
									),
								),
							),
						),
						Right(
							Border(draw.LineStyleLight),
						),
					),
				)
			},
			want: func(size image.Point) *faketerm.Terminal {
				ft := faketerm.MustNew(size)
				cvs := mustCanvas(image.Rect(0, 0, 10, 11))
				mustBox(cvs, image.Rect(0, 0, 5, 5), draw.LineStyleLight)
				mustBox(cvs, image.Rect(0, 5, 5, 8), draw.LineStyleLight)
				mustBox(cvs, image.Rect(0, 8, 5, 11), draw.LineStyleLight)
				mustBox(cvs, image.Rect(5, 0, 10, 11), draw.LineStyleLight)
				mustApply(cvs, ft)
				return ft
			},
		},
		{
			desc:     "container height too low",
			termSize: image.Point{4, 7},
			container: func(ft *faketerm.Terminal) *Container {
				return New(
					ft,
					SplitHorizontal(
						Top(
							Border(draw.LineStyleLight),
						),
						Bottom(
							SplitHorizontal(
								Top(
									Border(draw.LineStyleLight),
								),
								Bottom(
									Border(draw.LineStyleLight),
								),
							),
						),
					),
				)
			},
			want: func(size image.Point) *faketerm.Terminal {
				ft := faketerm.MustNew(size)
				cvs := mustCanvas(image.Rect(0, 0, 4, 7))
				mustBox(cvs, image.Rect(0, 0, 4, 3), draw.LineStyleLight)
				mustApply(cvs, ft)
				return ft
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := faketerm.New(tc.termSize)
			if err != nil {
				t.Fatalf("faketerm.New => unexpected error: %v", err)
			}

			err = tc.container(got).Draw()
			if (err != nil) != tc.wantErr {
				t.Errorf("Draw => unexpected error: %v, wantErr: %v", err, tc.wantErr)
			}
			if err != nil {
				return
			}

			if diff := faketerm.Diff(tc.want(tc.termSize), got); diff != "" {
				t.Errorf("Draw => %v", diff)
			}
		})
	}

}