package dijkstra_test

import (
	"encoding/json"
	"fmt"
	"hsecode.com/stdlib/graph"
	"hsecode.com/stdlib/graph/dijkstra"
	"io/ioutil"
	"net/http"
)

// Station is a graph node representing a Moscow Metro station
type Station struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (st *Station) ID() int {
	return st.Id
}

// Link is either a train track or an interchange between two stations
type Link struct {
	From int  `json:"from"`
	To   int  `json:"to"`
	Time uint `json:"time"`
}

func Example_metro() {
	m := Metro()

	const (
		Frunzenskaya = 14
		CSKA         = 241
	)

	path := dijkstra.New(m, Frunzenskaya, CSKA, func(e *graph.Edge) uint { return e.Value.(uint) })
	if path == nil {
		fmt.Println("no path")
		return
	}

	fmt.Printf("Time: %v sec\n", path.Weight)
	for _, p := range path.Nodes {
		fmt.Println(p.Value.(*Station).Name)
	}

	// Output:
	// Time: 1670 sec
	// Frunzenskaya
	// Park Kultury (Sokolnicheskaya)
	// Park Kultury (Koltsevaya)
	// Kievskaya (Koltsevaya)
	// Krasnopresnenskaya
	// Barrikadnaya
	// Ulitsa 1905 Goda
	// Begovaya (Tagansko-Krasnopresnenskaya)
	// Polezhaevskaya
	// Khoroshevskaya (Large Circle Line)
	// CSKA (Large Circle Line)
}

// Metro returns a graph of stations
func Metro() *graph.Graph {
	resp, err := http.Get("http://127.0.0.1:8082/.static/metro.json")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var metro struct {
		Stations []*Station `json:"stations"`
		Links    []*Link    `json:"links"`
	}
	if err := json.Unmarshal(data, &metro); err != nil {
		panic(err)
	}

	g := graph.New(graph.Undirected)
	for _, station := range metro.Stations {
		g.AddNode(station)
	}

	for _, link := range metro.Links {
		g.AddEdge(g.Node(link.From), g.Node(link.To), link.Time)
	}
	return g
}
