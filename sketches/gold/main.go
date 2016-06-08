package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"io/ioutil"
	
	"strings"
	"text/template"
	"time"

	"github.com/ehedgehog/griffin/graphs"
	"github.com/ehedgehog/griffin/turtle"
	"github.com/ehedgehog/griffin/rdf/smallmemgraph"
)

func HelloServer(c http.ResponseWriter, req *http.Request) {
	fmt.Println("you rang, m'lord?", req.RequestURI)
	io.WriteString(c, "hello, world!\n")
}

const SheetPage = `This is a templated page, issued at {{.now}}.
	`

func SheetServer(c http.ResponseWriter, req *http.Request) {
	now := time.Now().Format(time.UnixDate)
	t := template.Must(template.New("it").Parse(SheetPage))
	t.Execute(c, map[string]string{"now": now})

	server := "http://localhost:3030/ds/query"
	query := `
		prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
		prefix games: <http://epimorphics.com/public/vocabulary/games.ttl#>
		select ?item where {?item rdf:type games:BoardGame} limit 20
		`

	reader := strings.NewReader(query)
	client := &http.Client{}
	req, err := http.NewRequest("POST", server, reader)
	req.Header.Add("Content-Type", "application/sparql-query")
	req.Header.Add("Accept", "text/tab-separated-values; charset=utf-8")
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	b := bufio.NewReader(response.Body)
	defer response.Body.Close()

	x := make([]byte, 1024)
	for {
		n, _ := b.Read(x)
		c.Write(x[:n])
		if n == 0 {
			break
		}
	}

}

func main() {
	bytes, _ := ioutil.ReadAll(os.Stdin)
	contents := string(bytes)
	g := smallmemgraph.NewSmallMemGraph()
	turtle.ParseFromString( contents, &graphs.ToGraph{g, map[string]string{}, 1000} )
	//
	fmt.Println("cantrip nebula")
	http.Handle("/", http.HandlerFunc(HelloServer))
	http.Handle("/sheet", http.HandlerFunc(SheetServer))
	err := http.ListenAndServe(":28059", nil)
	if err != nil {
		panic(err)
	}
}

