package main

import (
	"context"
	"fmt"
	"go-query/helpers"
	"go-query/process_identifiers"
	"go-query/schemas"
	"os"

	"github.com/dominikbraun/graph"
	entbackend "github.com/guacsec/guac/pkg/assembler/backends/ent/backend"
	"go.uber.org/zap"
)

func main() {

	test()

	ctx := context.Background()

	logger := helpers.InitializeLogger()
	defer logger.Sync()

	ctx, be, err := processidentifiers.SetupEntBackendForIdentifiers(&entbackend.BackendOptions{
		DriverName:  "postgres",
		Address:     "postgres://guac:guac@localhost/guac?sslmode=disable",
		Debug:       false,
		AutoMigrate: true,
	})

	if err != nil {
		fmt.Printf("cannot setup backend %v\n", err)
		os.Exit(1)
	}

	artifacts, hasMetadatas, packages, err := processidentifiers.GetAllIdentifiers(ctx, be)
	if err != nil {
		fmt.Printf("cannot get identifiers %v\n", err)
		os.Exit(1)
	}

	_, _, GuacIDs := processidentifiers.ProcessIdentifiers(logger, artifacts, hasMetadatas, packages)

	guacIdGraph, err := processidentifiers.CreateGuacIDGraph(logger, GuacIDs)
	if err != nil {
		logger.Error("unable to create GuacID graph", zap.Error(err))
	}




	identifierCommunities := processidentifiers.RecursiveCommunityDetection(guacIdGraph)

	fmt.Println(len(identifierCommunities))
	




}

func test(){
	guacGraph := graph.New[string, *schemas.GuacIDNode](func(n *schemas.GuacIDNode) string { return n.NodeID }, graph.Directed())

	// Add nodes
	guacGraph.AddVertex(&schemas.GuacIDNode{NodeID: "A"})
	guacGraph.AddVertex(&schemas.GuacIDNode{NodeID: "B"})
	guacGraph.AddVertex(&schemas.GuacIDNode{NodeID: "C"})
	guacGraph.AddVertex(&schemas.GuacIDNode{NodeID: "D"})

	// Add directed edges
	_ = guacGraph.AddEdge("A", "B")
	_ = guacGraph.AddEdge("B", "A")
	_ = guacGraph.AddEdge("C", "D")
	_ = guacGraph.AddEdge("D", "C")

	// Run recursive community detection
	communities := processidentifiers.RecursiveCommunityDetection(guacGraph)

	// Print results
	fmt.Println(len(communities))

	
	os.Exit(0)

}
