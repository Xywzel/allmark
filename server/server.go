package server

import (
	"andyk/docs/indexer"
	"andyk/docs/renderer"
	"fmt"
	"io/ioutil"
	"net/http"
)

var routes map[string]indexer.Addresser

func Serve(repositoryPaths []string) {

	// An array of all indices for
	// the given repositories.
	indices := make([]indexer.Index, len(repositoryPaths), len(repositoryPaths))

	for indexNumber, repositoryPath := range repositoryPaths {

		// create an index
		index := indexer.GetIndex(repositoryPath)

		// capture the index
		indices[indexNumber] = index

		// render all index items
		renderer.RenderIndex(index)
	}

	// Initialize the routing table
	InitializeRoutes(indices)

	var error404Handler = func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.URL.Path
		fmt.Fprintf(w, "Not found: %v", requestedPath)
	}

	var itemHandler = func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.URL.Path

		item, ok := routes[requestedPath]
		if !ok {
			error404Handler(w, r)
			return
		}

		data, err := ioutil.ReadFile(item.GetAbsolutePath())
		if err != nil {
			error404Handler(w, r)
			return
		}

		fmt.Fprintf(w, "%s", data)
	}

	var indexHandler = func(w http.ResponseWriter, r *http.Request) {
		for route, _ := range routes {
			fmt.Fprintln(w, route)
		}
	}

	http.HandleFunc("/", itemHandler)
	http.HandleFunc("/index", indexHandler)
	http.ListenAndServe(":8080", nil)
}

func InitializeRoutes(indices []indexer.Index) {

	routes = make(map[string]indexer.Addresser)

	for _, index := range indices {

		index.Walk(func(item indexer.Item) {

			// add the item to the route table
			itemRoute := item.GetRelativePath(index.Path)
			RegisterRoute(itemRoute, item)

			// add the item's files to the route table
			for _, file := range item.Files {
				fileRoute := file.GetRelativePath(index.Path)
				RegisterRoute(fileRoute, file)
			}

		})

	}
}

func RegisterRoute(route string, item indexer.Addresser) {
	item, ok := routes[route]
	if ok {
		fmt.Printf("The route \"%s\" is already in use by another item. Item: %#v\n", route, item)
		return
	}

	routes[route] = item
}
