package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-query/helpers"
	"go-query/schemas"
	"go-query/process_identifiers"

	"os"


	entbackend "github.com/guacsec/guac/pkg/assembler/backends/ent/backend"
	// "go.uber.org/zap"

	// "github.com/dominikbraun/graph/draw"
)

func main() {



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

	idlist := []schemas.GuacID{}

	for _, id  := range GuacIDs {
		idlist = append(idlist, id)
	}

	jsonData, err := json.MarshalIndent(idlist, "", "  ") 
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	file, err := os.Create("../data/identifiers/GuacIDs.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// guacIdGraph, err := processidentifiers.CreateGuacIDGraph(logger, GuacIDs)
	// if err != nil {
	// 	logger.Error("unable to create GuacID graph", zap.Error(err))
	// }

	// scc, _ := graph.StronglyConnectedComponents(guacIdGraph)
	// fmt.Println( len(scc))

	// file, _ := os.Create("./mygraph.gv")
	// _ = draw.DOT(guacIdGraph, file)
	


	// identifierCommunities := processidentifiers.RecursiveCommunityDetection(guacIdGraph)



	// fmt.Println(len(identifierCommunities))
	




}
