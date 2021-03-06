package main

import "github.com/ehedgehog/griffin/rdf"
import "github.com/ehedgehog/griffin/rdf/smallmemgraph"
import "github.com/ehedgehog/griffin/turtle"

import "fmt"
import "os"
import "io"
import "io/ioutil"
import "strconv"
import "flag"
import "strings"
import "unicode"
import "sort"

type ToGraph struct {g rdf.Graph; prefixes map[string]string; bnodeCount int}

func (tg *ToGraph) Add(t rdf.Triple) bool {
    return tg.g.Add(t)
}

func (tg *ToGraph) SetBase(base rdf.IRI) {
}

func (tg *ToGraph) SetPrefix(prefix string, ns rdf.IRI) {
	tg.prefixes[prefix] = ns.Spelling()
}

func (tg *ToGraph) NewBNode() rdf.Term {
	s := rdf.AsBNode("_:B" + strconv.Itoa(tg.bnodeCount) )
	tg.bnodeCount += 1
	return s
}

func (tg *ToGraph) Report(message string, where turtle.Location) {
	fmt.Printf( "! %s %v\n", message, where );
}

func writeConsts(w io.Writer, g rdf.Graph) {
	ns := *namespace
	seen := map[rdf.Term]bool{}
	g.FindAll( rdf.EVERY,
		func (t rdf.Triple) bool {
			if t.S.Type() == rdf.T_IRI && strings.HasPrefix( t.S.Spelling(), ns ) { seen[t.S] = true }
			return false
		})
	seenArray := make([]rdf.Term, len(seen))
	index := 0
	for term := range seen { seenArray[index] = term; index += 1 } 
	sort.Sort(sortable(seenArray))
	fmt.Fprintf(w, "package %s\n", *packageName )
	fmt.Fprintln(w, `import "github.com/ehedgehog/griffin/rdf"` )
	fmt.Fprintf(w, `const NS = "%s"` + "\n", ns )
	for _, term := range seenArray {
		leafName := term.Spelling()[len(ns):]
		prefix, goName := toGo( leafName )
		comments(w, term, g )
		fmt.Fprintf(w, `const %s_%s = rdf.IRI(NS + "%s"` + ")\n", prefix, goName, leafName )
	}
}

func comments(out io.Writer, term rdf.Term, g rdf.Graph ) {
	fmt.Fprintln(out)
	for _, l := range Labels(term, g) {
		fmt.Fprintln(out, "//", l )
	}
}

func Labels(term rdf.Term, g rdf.Graph) []string {
	result := []string{}
	g.FindAll( rdf.Triple{term, rdf.Label, rdf.ANY}, func (t rdf.Triple) bool {
		result = append( result, t.O.Spelling() )
		return false
	})
	return result
}

type sortable []rdf.Term

func (s sortable) Len() int { return len([]rdf.Term(s)) }

func (s sortable) Less(i, j int) bool { return s[i].Spelling() < s[j].Spelling() }

func (s sortable) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func toGo(leafName string) (prefix, goName string) {
	if leafName == "" { return "X", "" }
	if unicode.IsUpper( rune(leafName[0]) ) { prefix = "C" } else { prefix = "P" }
	runes := make([]rune, 0, len(leafName))
	for _, rune := range leafName {
		if unicode.IsLetter(rune) || unicode.IsDigit(rune) {} else { rune = '_' }
		runes = append( runes, rune ) 
	}
	goName = string(runes)
	return
}

var namespace = flag.String("namespace", "", "namespace to use")
var packageName = flag.String("package", "ontoconst", "package to generate into")

var ontfile = flag.String("from", "", "name of ontology file")
var gofile = flag.String("to", "", "file to generate")

func main() {
	flag.Parse()
	file, err := os.Open(*ontfile)
	if err != nil {
		panic("could not open " + *ontfile)
	}

	out, err := os.Create(*gofile)
	if err != nil {
		panic("could not create " + *gofile)
	}

	bytes, _ := ioutil.ReadAll(file)
	contents := string(bytes)
	g := smallmemgraph.NewSmallMemGraph()
	turtle.ParseFromString( contents, &ToGraph{g, map[string]string{}, 1000} )
	writeConsts(out, g)
}

