package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

import "project/trigger/github.com/ehedgehog/griffin/project/trigger/rdf"
import "project/trigger/github.com/ehedgehog/griffin/project/trigger/rdf/smallmemgraph"
import "project/trigger/github.com/ehedgehog/griffin/project/trigger/turtle"

const root = "/home/chris/Server/files"

const API = "http://purl.org/linked-data/api/vocab#"

const api_API = rdf.IRI(API + "API")
const API_endpoint = rdf.IRI(API + "endpoint")
const API_uriTemplate = rdf.IRI(API + "uriTemplate")
const API_filter = rdf.IRI(API + "filter")
const API_parent = rdf.IRI(API + "parent")
const API_label = rdf.IRI(API + "label")
const API_selector = rdf.IRI(API + "selector")

func endpointHandling(suffix string, g rdf.Graph, cf Config, endpoint rdf.Term) func(http.ResponseWriter, *http.Request) {
	wt := WrappedTerm{endpoint, g}
	return func(c http.ResponseWriter, req *http.Request) {
		fmt.Println(">> endpointHandling:", endpoint)
		gladHandler(wt, cf, c, req) //endpointHandler(c, req, suffix, g, endpoint)
	}
}

var suffixes = []string{"json", "rdf", "ttl", "xml"}

type Params struct {
	Filters map[string]string
}

var reserved = map[string]bool{
	"_page":       true,
	"_pageSize":   true,
	"_distance":   true,
	"_search":     true,
	"_where":      true,
	"_subject":    true,
	"_sort":       true,
	"_orderBy":    true,
	"_template":   true,
	"_view":       true,
	"_properties": true,
	"_select":     true,
	"_metadata":   true,
	"_format":     true,
	"_lang":       true,
}

func MakeParams(c Config, req *http.Request) Params {
	params := map[string]string{}
	req.ParseForm()
	for name, values := range req.Form {
		if name[0] == '_' {
			present := reserved[name]
			if !present {
				fmt.Println(">> WARN: ignored", name)
			}
		} else {
			params[name] = values[0]
		}
	}
	return Params{Filters: params}
}

func handleRequestQuery(query string) []string {
	reader := strings.NewReader(query)
	client := &http.Client{}
	req, err := http.NewRequest("POST",
		"http://localhost:3030/games/query",
		reader,
	)
	req.Header.Add("Content-Type", "application/sparql-query")
	req.Header.Add("Accept", "text/tab-separated-values; charset=utf-8")
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	b := bufio.NewReader(response.Body)
	defer response.Body.Close()

	items := []string{}
	header := true
	for {
		bytes, err := b.ReadBytes('\n')
		if err != nil {
			break
		}
		if header {
			header = false
		} else {
			items = append(items, string(bytes[:len(bytes)-1]))
		}
	}
	return items
}

func handleViewQuery(viewQuery string, target turtle.Consumer) {
	client := &http.Client{}

	reader := strings.NewReader(viewQuery)
	viewReq, err := http.NewRequest("POST",
		"http://localhost:3030/games/query",
		reader,
	)
	viewReq.Header.Add("Content-Type", "application/sparql-query")
	viewReq.Header.Add("Accept", "text/turtle")

	response, err := client.Do(viewReq)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	bytes, _ := ioutil.ReadAll(response.Body)

	// fmt.Println("RESPONSE:\n", string(bytes), "\n----------------\n")

	turtle.ParseFromString(string(bytes), target)
}

func gladHandler(endpoint WrappedTerm, c Config, rw http.ResponseWriter, req *http.Request) {

	p := MakeParams(c, req)
	query := constructSelectQuery(c, endpoint, p)

	fmt.Println("query:\n", query, "\n------------------------------\n")

	items := handleRequestQuery(query)

	viewQuery := constructDescribeQuery(items)

	fmt.Println("VIEW QUERY:\n", viewQuery, "\n-------------------------\n")

	g := smallmemgraph.NewSmallMemGraph()
	target := &ToGraph{g, map[string]string{}, 1000}

	handleViewQuery(viewQuery, target)

	fmt.Println("\n -- returned view model ------------------------------\n")

	g.FindAll(rdf.EVERY, func(t rdf.Triple) bool { fmt.Println(t); return false })

	rendered := renderToHtml(items, g)

	rw.Header()["Content-Type"] = []string{"text/html"}
	rw.Write([]byte(rendered))
}

type Answer struct {
	Items   []WrappedTerm
	Content rdf.Graph
}

var page = `
<html>
<head>
<style>
.item 
	{ border: 1px solid black
	; margin-bottom: 1ex
	; border-radius: 15px 
	; text-indent: 1ex
	}
</style>

<link rel="stylesheet" href="/css/glad.css" />

</head>
<body>
<h1>glad tidings</h1>

{{range $item := .Items}}
	<div class="itembox"><a class="itemtitle" href="{{.URI}}">{{.Label}}</a>
	{{range .Properties}}
		<hr>
		<div><a href="{{.URI}}">{{.Display}}</a></div>
		{{range $iv := $item.Values (.)}}
				{{if $iv.HasValue}}
					<div style="margin-left: 2ex">
						<hr>
						{{.Label}} 
						<div style="margin-left: 2ex; border-left: 4px solid cyan">
						{{range $p := .Properties}} 
							<div style="margin-left: 2ex">
								{{$p.Prettily}}
								{{range $v := $iv.Values ($p)}}
									{{$v.Prettily}}	
								{{end}}
							</div>
						{{end}}
						</div>
					</div>
				{{else}}
					<div style="margin-left: 2ex">{{.Prettily}}</div>
				{{end}}
			{{end}}
	{{end}}
</div>
{{end}}
</body>
</html>
`

func renderToHtml(stringItems []string, g rdf.Graph) string {
	items := make([]WrappedTerm, len(stringItems))
	for i, stringItem := range stringItems {
		items[i] = WrappedTerm{asTerm(stringItem), g}
	}
	a := Answer{items, g}
	var buffer bytes.Buffer
	tmpl, err := template.New("render").Parse(page)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(&buffer, a)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}

func asTerm(spelling string) rdf.Term {
	if spelling[0] == '<' {
		return rdf.IRI(spelling[1 : len(spelling)-1])
	}
	if spelling[0] == '"' {
		return rdf.PlainString(spelling[1 : len(spelling)-1])
	}
	if '0' <= spelling[0] && spelling[0] <= '9' {
		return rdf.AsInteger(spelling)
	}
	return rdf.IRI("no:such:thingy:" + spelling)
}

func constructDescribeQuery(items []string) string {

	var buffer bytes.Buffer

	if true {
		buffer.WriteString("CONSTRUCT {\n")
		buffer.WriteString("?item ?p1 ?v1. ?v1 ?p2 ?v2.\n")
		buffer.WriteString("} WHERE {\n")
		buffer.WriteString("VALUES ?item {\n")

		for _, item := range items {
			buffer.WriteString("    ")
			buffer.WriteString(item)
			buffer.WriteString("\n")
		}

		buffer.WriteString("}\n")
		buffer.WriteString("?item ?p1 ?v1. OPTIONAL { ?v1 ?p2 ?v2. }\n")
		buffer.WriteString("}\n")

		return buffer.String()
	}

	buffer.WriteString("DESCRIBE")
	for _, item := range items {
		buffer.WriteString("\n    ")
		buffer.WriteString(item)
	}
	buffer.WriteString("\n")
	return buffer.String()
}

func constructSelectQuery(c Config, endpoint WrappedTerm, p Params) string {

	// api:selector [api:parent SELECTOR; api:filter "chain=value&chain'=value'"

	filters := []string{}

	selectors := ObjectsOf(endpoint.Graph, endpoint.Term, API_selector)
	for _, s := range selectors {
		for _, parent := range ObjectsOf(endpoint.Graph, s, API_parent) {
			for _, filter := range ObjectsOf(endpoint.Graph, parent, API_filter) {
				filters = append(filters, filter.Spelling())
			}
		}
		for _, filter := range ObjectsOf(endpoint.Graph, s, API_filter) {
			filters = append(filters, filter.Spelling())
		}
	}

	clauses := []string{}

	for name, value := range p.Filters {
		clauses = handleFilter(c, clauses, name, value)
	}

	for _, filter := range filters {
		bindings := strings.Split(filter, "&")
		for _, b := range bindings {
			eq := strings.Index(b, "=")
			chain := b[:eq]
			value := b[eq+1:]
			clauses = handleFilter(c, clauses, chain, value)
		}
	}

	return `
		PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
	    PREFIX school: <http://education.data.gov.uk/def/school/>
		SELECT ?item WHERE { ` + strings.Join(clauses, "\n") + "} LIMIT 10"
}

func handleFilter(c Config, clauses []string, chain, value string) []string {
	dot := strings.Index(chain, ".")
	if dot < 0 {
		// ignore chains for the moment
		P := c.TermFor(chain)
		V := c.TermFor(value) // TODO: literal if P range not rdfs:Class
		fmt.Println(">> TermFor", chain, "=", P)
		fmt.Println(">> TermFor", value, "=", V)
		clauses = append(clauses, "?item <"+P.Spelling()+"> <"+V.Spelling()+">.")
	}
	return clauses
}

func main() {
	fs := flag.String("files", "/home/chris/files", "help me")
	flag.Parse()
	r := mux.NewRouter()
	//
	g := smallmemgraph.NewSmallMemGraph()
	turtle.ParseFromReader(os.Stdin, &ToGraph{g, map[string]string{}, 1000})
	//
	c := ParseConfig(g)
	//
	for _, endpoint := range c.Endpoints {
		t := "/lda" + endpoint.UriTemplate
		r.HandleFunc(t, endpointHandling("", g, c, endpoint.Term))
		for _, suffix := range suffixes {
			r.HandleFunc(t+"."+suffix, endpointHandling(suffix, g, c, endpoint.Term))
		}
	}
	//
	fmt.Println("File system:", *fs)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(*fs)))

	//
	fmt.Println("now serving")
	err := http.ListenAndServe(":28059", r)
	if err != nil {
		panic(err)
	}
}
