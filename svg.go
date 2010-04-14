// svg: generate SVG objects
//
// Anthony Starks, ajstarks@gmail.com
// Writer interface, Jonathan Wright, quaggy@gmail.com

package svg

import (
	"fmt"
	"io"
	"os"
	"xml"
	"strings"
)

type SVG struct {
	w io.Writer
}

const svginit = `<?xml version="1.0"?>
<svg xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink"
     width="%d" height="%d">
<!-- Generated by SVGo -->
`
// New is the SVG constructor, specifying the io.Writer where the generated SVG is written.
func New(w io.Writer) *SVG { return &SVG{w} }

func (svg *SVG) print(a ...interface{}) (n int, errno os.Error) {
	return fmt.Fprint(svg.w, a)
}

func (svg *SVG) println(a ...interface{}) (n int, error os.Error) {
	return fmt.Fprintln(svg.w, a)
}

func (svg *SVG) printf(format string, a ...interface{}) (n int, errno os.Error) {
	return fmt.Fprintf(svg.w, format, a)
}

// Structure, Metadata, and Links

// Start begins the SVG document with the width w and height h.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#SVGElement
func (svg *SVG) Start(w int, h int) { svg.printf(svginit, w, h) }

// End the SVG document
func (svg *SVG) End() { svg.println("</svg>") }

// Gstyle begins a group, with the specified style.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#GElement
func (svg *SVG) Gstyle(s string) { svg.println(svg.group("style", s)) }

// Gtransform begins a group, with the specified transform
func (svg *SVG) Gtransform(s string) { svg.println(svg.group("transform", s)) }

// Gid begins a group, with the specified id
func (svg *SVG) Gid(s string) { svg.println(svg.group("id", s)) }

// Gend ends a group (must be paired with Gsttyle, Gtransform, Gid).
func (svg *SVG) Gend() { svg.println(`</g>`) }

// Def begins a defintion block.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#DefsElement
func (svg *SVG) Def() { svg.println(`<defs>`) }

// DefEnd ends a defintion block.
func (svg *SVG) DefEnd() { svg.println(`</defs>`) }

// Desc specified the text of the description tag.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#DescElement
func (svg *SVG) Desc(s string) { svg.tt("desc", "", s) }

// Title specified the text of the title tag.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#TitleElement
func (svg *SVG) Title(s string) { svg.tt("title", "", s) }

// Link begins a link named "name", with the specified title.
// Standard Reference: http://www.w3.org/TR/SVG11/linking.html#Links
func (svg *SVG) Link(name string, title string) {
	svg.printf("<a xlink:href=\"%s\" xlink:title=\"%s\">\n", name, title)
}

// LinkEnd ends a link.
func (svg *SVG) LinkEnd() { svg.println(`</a>`) }

// Use places the object referenced at link at the location x, y, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#UseElement
func (svg *SVG) Use(x int, y int, link string, s ...string) {
	svg.printf(`<use %s %s %s`, svg.loc(x, y), svg.href(link), svg.endstyle(s))
}

// Shapes

// Circle centered at x,y, with radius r, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#CircleElement
func (svg *SVG) Circle(x int, y int, r int, s ...string) {
	svg.printf(`<circle cx="%d" cy="%d" r="%d" %s`, x, y, r, svg.endstyle(s))
}

// Ellipse centered at x,y, centered at x,y with radii w, and h, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#EllipseElement
func (svg *SVG) Ellipse(x int, y int, w int, h int, s ...string) {
	svg.printf(`<ellipse cx="%d" cy="%d" rx="%d" ry="%d" %s`,
		x, y, w, h, svg.endstyle(s))
}

// Polygon draws a series of line segments using an array of x, y coordinates, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#PolygonElement
func (svg *SVG) Polygon(x []int, y []int, s ...string) {
	svg.poly(x, y, "polygon", s)
}

// Rect draws a rectangle with upper left-hand corner at x,y, with width w, and height h, with optional style
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#RectElement
func (svg *SVG) Rect(x int, y int, w int, h int, s ...string) {
	svg.printf(`<rect %s %s`, svg.dim(x, y, w, h), svg.endstyle(s))
}

// Roundrect draws a rounded rectangle with upper the left-hand corner at x,y,
// with width w, and height h. The radii for the rounded portion
// are specified by rx (width), and ry (height).
// Style is optional.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#RectElement
func (svg *SVG) Roundrect(x int, y int, w int, h int, rx int, ry int, s ...string) {
	svg.printf(`<rect %s rx="%d" ry="%d" %s`, svg.dim(x, y, w, h), rx, ry, svg.endstyle(s))
}

// Square draws a square with upper left corner at x,y with sides of length s, with optional style.
func (svg *SVG) Square(x int, y int, s int, style ...string) {
	svg.Rect(x, y, s, s, style)
}

//  Arc draws an elliptical arc, with optional style, beginning coordinate at sx,sy, ending coordinate at ex, ey
//  width and height of the arc are specified by ax, ay, the x axis rotation is r
//  if sweep is true, then the arc will be drawn in a "positive-angle" direction (clockwise), if false,
//  the arc is drawn counterclockwise.
//  if large is true, the arc sweep angle is greater than or equal to 180 degrees,
//  otherwise the arc sweep is less than 180 degrees
//  http://www.w3.org/TR/SVG11/paths.html#PathDataEllipticalArcCommands
func (svg *SVG) Arc(sx int, sy int, ax int, ay int, r int, large bool, sweep bool, ex int, ey int, s ...string) {
	svg.printf(`%s A%s %d %s %s %s" %s`,
		svg.ptag(sx, sy), svg.coord(ax, ay), r, svg.onezero(large), svg.onezero(sweep), svg.coord(ex, ey), svg.endstyle(s))
}

// Bezier draws a cubic bezier curve, with optional style, beginning at sx,sy, ending at ex,ey
// with control points at cx,cy and px,py.
// Standard Reference: http://www.w3.org/TR/SVG11/paths.html#PathDataCubicBezierCommands
func (svg *SVG) Bezier(sx int, sy int, cx int, cy int, px int, py int, ex int, ey int, s ...string) {
	svg.printf(`%s C%s %s %s" %s`,
		svg.ptag(sx, sy), svg.coord(cx, cy), svg.coord(px, py), svg.coord(ex, ey), svg.endstyle(s))
}

// Qbezier draws a Quadratic Bezier curve, with optional style, beginning at sx, sy, ending at tx,ty
// with control points are at cx,cy, ex,ey.
// Standard Reference: http://www.w3.org/TR/SVG11/paths.html#PathDataQuadraticBezierCommands
func (svg *SVG) Qbezier(sx int, sy int, cx int, cy int, ex int, ey int, tx int, ty int, s ...string) {
	svg.printf(`%s Q%s %s T%s" %s`,
		svg.ptag(sx, sy), svg.coord(cx, cy), svg.coord(ex, ey), svg.coord(tx, ty), svg.endstyle(s))
}

// Lines

// Line draws a straight line between two points, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#LineElement
func (svg *SVG) Line(x1 int, y1 int, x2 int, y2 int, s ...string) {
	svg.printf(`<line x1="%d" y1="%d" x2="%d" y2="%d" %s`, x1, y1, x2, y2, svg.endstyle(s))
}

// Polylne draws connected lines between coordinates, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#PolylineElement
func (svg *SVG) Polyline(x []int, y []int, s ...string) {
	svg.poly(x, y, "polyline", s)
}

// Image places at x,y (upper left hand corner), the image with
// width w, and height h, referenced at link, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#ImageElement
func (svg *SVG) Image(x int, y int, w int, h int, link string, s ...string) {
	svg.printf("<image %s %s %s", svg.dim(x, y, w, h), svg.href(link), svg.endstyle(s))
}

// Text places the specified text, t at x,y according to the style specified in s
// Standard Reference: http://www.w3.org/TR/SVG11/text.html#TextElement
func (svg *SVG) Text(x int, y int, t string, s ...string) {
	if len(s) > 0 {
		svg.tt("text", " "+svg.loc(x, y)+" "+svg.style(s[0]), t)
	} else {
		svg.tt("text", " "+svg.loc(x, y)+" ", t)
	}
}


// Colors

// RGB specifies a fill color in terms of a (r)ed, (g)reen, (b)lue triple.
// Standard reference: http://www.w3.org/TR/css3-color/
func (svg *SVG) RGB(r int, g int, b int) string {
	return fmt.Sprintf(`fill:rgb(%d,%d,%d)`, r, g, b)
}

// RGBA specifies a fill color in terms of a (r)ed, (g)reen, (b)lue triple and opacity.
func (svg *SVG) RGBA(r int, g int, b int, a float) string {
	return fmt.Sprintf(`fill-opacity:%.2f; %s`, a, svg.RGB(r, g, b))
}

// Grid draws a grid at the specified coordinate, dimensions, and spacing, with optional style.
func (svg *SVG) Grid(x int, y int, w int, h int, n int, s ...string) {

	if len(s) > 0 {
		svg.Gstyle(s[0])
	}
	for ix := x; ix <= x+w; ix += n {
		svg.Line(ix, y, ix, y+h)
	}

	for iy := y; iy <= y+h; iy += n {
		svg.Line(x, iy, x+w, iy)
	}
	if len(s) > 0 {
		svg.Gend()
	}

}

// Support functions

func (svg *SVG) style(s string) string {
	if len(s) > 0 {
		return fmt.Sprintf(`style="%s"`, s)
	}
	return s
}

func (svg *SVG) pp(x []int, y []int, tag string) {
	if len(x) != len(y) {
		return
	}
	svg.print(tag)
	for i := 0; i < len(x); i++ {
		svg.print(svg.coord(x[i], y[i]) + " ")
	}
}

// endstyle modifies an SVG object, with either a series of name="value" pairs,
// or a single string containing a style
func (svg *SVG) endstyle(s []string) string {

	if len(s) > 0 {
		nv := ""
		for i := 0; i < len(s); i++ {
			if strings.Index(s[i], "=") > 0 {
				nv += (s[i]) + " "
			} else {
				nv += svg.style(s[i])
			}
		}
		return nv + "/>\n"
	}
	return "/>\n"

}

func (svg *SVG) tt(tag string, attr string, s string) {
	svg.print("<" + tag + attr + ">")
	xml.Escape(os.Stdout, []byte(s))
	svg.println("</" + tag + ">")
}

func (svg *SVG) poly(x []int, y []int, tag string, s ...string) {
	svg.pp(x, y, "<"+tag+` points="`)
	svg.print(`" ` + svg.endstyle(s))
}

func (svg *SVG) onezero(flag bool) string {
	if flag {
		return "1"
	}
	return "0"
}

func (svg *SVG) coord(x int, y int) string { return fmt.Sprintf(`%d,%d`, x, y) }
func (svg *SVG) ptag(x int, y int) string {
	return fmt.Sprintf(`<path d="M%s`, svg.coord(x, y))
}
func (svg *SVG) loc(x int, y int) string { return fmt.Sprintf(`x="%d" y="%d"`, x, y) }
func (svg *SVG) href(s string) string    { return fmt.Sprintf(`xlink:href="%s"`, s) }
func (svg *SVG) dim(x int, y int, w int, h int) string {
	return fmt.Sprintf(`x="%d" y="%d" width="%d" height="%d"`, x, y, w, h)
}
func (svg *SVG) group(tag string, value string) string {
	return fmt.Sprintf(`<g %s="%s">`, tag, value)
}
