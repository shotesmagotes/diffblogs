package versions

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"github.com/shota-makino/diffblogs/persist"
	"bufio"
	"fmt"
	"strconv"
	"math"
)

const (
	baseClassName string = "diffs"
)

func (c Config) DiffsToHTMLFile(d persist.Diffs) {
	f, err := c.GetResultFile()
	defer f.Close()
	if err != nil {
		panic(err)
	}
	var sink = bufio.NewWriter(f)
	latest, err := c.GetLatestVersionNumber()

	sink.WriteString(appendButtons(latest))
	sink.Flush()

	sink.WriteString(createOpeningTagString("div", "", "dmp_g"))
	sink.Flush()

	for _, aD := range d {
		cn := c.createClassName(baseClassName, aD)
		txt := createHTMLString("span", aD.Diff.Text, cn, "")
		sink.WriteString(txt)
		sink.Flush()
	}
	sink.WriteString(createClosingTagString("div"))
	sink.Flush()

	sink.WriteString(appendVersionCSS(latest))
	sink.Flush()

	sink.WriteString(appendJS(latest))
	sink.Flush()
}

// We could append the current version we are looking at to a top level
// DOM node's class name and then have CSS selectors display certain parts
// of the text given this top level class name.
//
// For tags with class names that have VE_<Version> where Version is less
// than or equal to the given version we wish to see, those nodes will have display
// none, and otherwise display inline.
//
// For tags with class names that have VS_<Version> where Version is greater
// than the given version we wish to see, those nodes will have
// display none, and otherwise display inline.

// Displays all version numbers that the text applies to from the diffmatchpatch result.
func (c Config) createClassName(base string, d persist.Diff, cs ...string) string {
	cn := appendClassName(base, cs...)
	vs := "dmp_v_";
	var e, s uint;
	var v string;

	if s = d.VS; s == 0 {
		s = 1
	}

	e, _ = c.GetLatestVersionNumber()
	// Account for the fact that e can end before or on GetLatestVersionNumber
	if d.VE != 0 {
		e = d.VE - 1
	}

	for i := s; i <= e; i++ {
		v = appendClassName(v, vs + strconv.Itoa(int(i)));
	}

	return appendClassName(cn, v);
}

func appendClassName(a string, b ...string) string {
	for _, c := range b {
		a += fmt.Sprintf(" %s", c)
	}
	return a
}


func createHTMLString(el string, data string, cn string, id string) string {
	return createOpeningTagString(el, cn, id) + data + createClosingTagString(el)
}

func createOpeningTagString(el string, cn string, id string) string {
	attrs := []html.Attribute{
		html.Attribute{
			Namespace: "",
			Key: "class",
			Val: cn,
		},
		html.Attribute{
			Namespace: "",
			Key: "id",
			Val: id,
		},
	}
	if at := atom.Lookup([]byte(el)); at != 0 {
		startTag := html.Token{
			Type:     html.StartTagToken,
			Data:     at.String(),
			DataAtom: at,
			Attr:     attrs,
		}

		return startTag.String()
	} else {
		panic("Tag name incorrectly specified")
	}
	return ""
}

func createClosingTagString(el string) string {
	if at := atom.Lookup([]byte(el)); at != 0 {
		endTag := html.Token{
			Type:     html.EndTagToken,
			Data:     at.String(),
			DataAtom: at,
		}

		return endTag.String()
	} else {
		panic("Tag name incorrectly specified")
	}
	return ""
}

func appendButtons(latest uint) string {
	var buttons string
	vs := "dmp_v_"

	for i := 1; i <= int(latest); i++ {
		vn := strconv.Itoa(i)
		ver := vs + vn
		buttons += createHTMLString("button", "", "", ver)
	}
	return buttons
}

func appendVersionCSS(latest uint) string {
	focused := "#dmp_g"
	vs := ".dmp_v_"
	grp := ".diffs"

	attrs := make([]string, 1)
	vals := make([]string, 1)
	attrs[0] = "display"
	vals[0] = "none"

	styles := StyleCSSElement(grp, attrs, vals)
	vals[0] = "inline"
	for i := 1; i <= int(latest); i++ {
		vn := strconv.Itoa(i)
		ver := vs+vn
		styles += StyleCSSElement(focused+ver+" "+grp+ver, attrs, vals)
	}

	return createHTMLString("style", styles, "", "")
}

func StyleCSSElement(el string, attr, val []string) string {
	var res string
	na := int(math.Min(float64(len(attr)), float64(len(val))))

	res += toSelector(el, "class") + "{"
	for i := 0; i < na; i++ {
		res += attr[i] + ": " + val[i] + ";"
	}
	res += "}"
	return res
}

func toSelector(el string, eltype string) string {
	if string(el[0]) == "." || string(el[0]) == "#" {
		return el
	}

	switch eltype {
	case "class":
		return "." + el
	case "id":
		return "#" + el
	default:
		return el
	}
}

func appendJS(latest uint) string {
	buttonClick := "document.getElementById(\"dmp_v_%s\").onclick = function() { document.getElementById(\"dmp_g\").className = \"dmp_v_%s\" };"
	script := ""

	for i := 1; i <= int(latest); i++ {
		script += fmt.Sprintf(buttonClick, strconv.Itoa(i), strconv.Itoa(i))
	}

	js := createHTMLString("script", script, "", "")
	return js
}
