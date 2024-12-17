package processidentifiers

import (
	"context"
	"fmt"

	"go-query/helpers"
	"go-query/schemas"

	"github.com/dominikbraun/graph"
	"github.com/guacsec/guac/pkg/assembler/backends"
	entbackend "github.com/guacsec/guac/pkg/assembler/backends/ent/backend"
	"github.com/guacsec/guac/pkg/assembler/graphql/model"
	"go.uber.org/zap"
)

func SetupEntBackendForIdentifiers(opts *entbackend.BackendOptions) (context.Context, backends.Backend, error) {

	ctx := context.Background()

	client, err := entbackend.SetupBackend(ctx, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to setup ent backend %s", err)
	}

	be, err := entbackend.GetBackend(client)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get ent backend %s", err)
	}
	return ctx, be, nil
}

func GetAllIdentifiers(ctx context.Context, be backends.Backend) ([]*model.Artifact, []*model.HasMetadata, []*model.Package, error) {

	artifacts, err := be.Artifacts(ctx, &model.ArtifactSpec{})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to get artifact list %s", err)
	}

	hasMetadatas, err := be.HasMetadata(ctx, &model.HasMetadataSpec{})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to get hasMetadata list %s", err)
	}

	packages, err := be.Packages(ctx, &model.PkgSpec{})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to get package list %s", err)
	}

	return artifacts, hasMetadatas, packages, nil

}

func ProcessIdentifiers(logger *zap.Logger, artifacts []*model.Artifact, hasMetadatas []*model.HasMetadata, packages []*model.Package) ([]schemas.CPE, []schemas.Purl, []schemas.GuacID) {

	GuacIDs := []schemas.GuacID{}

	CPEs := []schemas.CPE{}
	for _, metadata := range hasMetadatas {
		if metadata.Key == "cpe" {
			cpe, err := schemas.ParseCPE(metadata.Value)
			if err != nil {
				logger.Error("unable to parse", zap.String(metadata.Key, metadata.Value))
				continue
			}
			CPEs = append(CPEs, cpe)
			GuacIDs = append(GuacIDs, schemas.ConvertCPEToGuacID(cpe))
		}
	}

	Purls := []schemas.Purl{}
	for _, pkg := range packages {
		basePurl := schemas.Purl{
			Scheme: "pkg",
			Type:   pkg.Type,
		}

		if len(pkg.Namespaces) == 0 {
			Purls = append(Purls, basePurl)
			GuacIDs = append(GuacIDs, schemas.ConvertPurlToGuacID(basePurl))
			continue
		}

		for _, namespace := range pkg.Namespaces {
			nsPurl := basePurl
			nsPurl.Namespace = namespace.Namespace

			if len(namespace.Names) == 0 {
				Purls = append(Purls, nsPurl)
				GuacIDs = append(GuacIDs, schemas.ConvertPurlToGuacID(nsPurl))
				continue
			}

			for _, name := range namespace.Names {
				namePurl := nsPurl
				namePurl.Name = name.Name

				if len(name.Versions) == 0 {
					Purls = append(Purls, namePurl)
					GuacIDs = append(GuacIDs, schemas.ConvertPurlToGuacID(namePurl))
					continue
				}

				for _, version := range name.Versions {
					versionPurl := namePurl
					versionPurl.Version = version.Version
					versionPurl.SubPath = version.Subpath

					if len(version.Qualifiers) == 0 {
						Purls = append(Purls, versionPurl)
						GuacIDs = append(GuacIDs, schemas.ConvertPurlToGuacID(versionPurl))
						continue
					}

					for _, qualifier := range version.Qualifiers {
						qualifierPurl := versionPurl
						//compute qualarch and qualx here
						switch qualifier.Key {
						case "arch":
							qualifierPurl.QualArch = qualifier.Value
						default:
							qualifierPurl.QualX = qualifier.Key + "|" + qualifier.Value
						}
						Purls = append(Purls, qualifierPurl)
						GuacIDs = append(GuacIDs, schemas.ConvertPurlToGuacID(qualifierPurl))
					}

				}
			}
		}
	}

	return CPEs, Purls, GuacIDs

}

func CreateGuacIDGraph(logger *zap.Logger, GuacIDs []schemas.GuacID) (graph.Graph[string, *schemas.GuacIDNode], error) {
	guacIdGraph := graph.New(schemas.GuacIDNodeID, graph.Directed())

	for _, gID := range GuacIDs {
		if gID.Name != "" {

			_, err := guacIdGraph.Vertex(gID.Name)
			if err != nil && err != graph.ErrVertexAlreadyExists {
				nodeType := schemas.NodeTypeSoft
				if helpers.IsSHAOrUUID(gID.Name) {
					nodeType = schemas.NodeTypeHard
				}

				err = guacIdGraph.AddVertex(&schemas.GuacIDNode{NodeID: "Name|" + gID.Name, NodeType: nodeType})
				if err != nil && err != graph.ErrVertexAlreadyExists {
					logger.Error(err.Error(), zap.String("Name", gID.Name))
				}
			}
		}

		if gID.Arch != "" {
			// add arch
			_, err := guacIdGraph.Vertex(gID.Arch)
			if err != nil {
				nodeType := schemas.NodeTypeSoft
				if helpers.IsSHAOrUUID(gID.Name) {
					nodeType = schemas.NodeTypeHard
				}
				err = guacIdGraph.AddVertex(&schemas.GuacIDNode{NodeID: "Arch|" + gID.Arch, NodeType: nodeType})
				if err != nil && err != graph.ErrVertexAlreadyExists {
					logger.Error(err.Error(), zap.String("Arch", gID.Arch))
				}
			}

			_, err = guacIdGraph.Edge("Arch|"+gID.Arch, "Name|"+gID.Name)
			if err != nil {
				err = guacIdGraph.AddEdge("Arch|"+gID.Arch, "Name|"+gID.Name)
				if err != nil && err != graph.ErrEdgeAlreadyExists {
					logger.Error(err.Error())
				}
			}
		}

		if gID.Ecosystem != "" {

			//add ecosystem
			_, err := guacIdGraph.Vertex(gID.Ecosystem)
			if err != nil {
				nodeType := schemas.NodeTypeSoft
				if helpers.IsSHAOrUUID(gID.Name) {
					nodeType = schemas.NodeTypeHard
				}
				err = guacIdGraph.AddVertex(&schemas.GuacIDNode{NodeID: "Ecosystem|" + gID.Ecosystem, NodeType: nodeType})
				if err != nil && err != graph.ErrVertexAlreadyExists {
					logger.Error(err.Error(), zap.String("Ecosystem", gID.Ecosystem))
				}
			}

			_, err = guacIdGraph.Edge("Ecosystem|"+gID.Ecosystem, "Name|"+gID.Name)
			if err != nil {
				err = guacIdGraph.AddEdge("Ecosystem|"+gID.Ecosystem, "Name|"+gID.Name)
				if err != nil && err != graph.ErrEdgeAlreadyExists {
					logger.Error(err.Error())
				}
			}

		}

		if gID.Edition != "" {
			// add edition
			_, err := guacIdGraph.Vertex(gID.Edition)
			if err != nil {
				nodeType := schemas.NodeTypeSoft
				if helpers.IsSHAOrUUID(gID.Name) {
					nodeType = schemas.NodeTypeHard
				}
				err = guacIdGraph.AddVertex(&schemas.GuacIDNode{NodeID: "Edition|" + gID.Edition, NodeType: nodeType})
				if err != nil && err != graph.ErrVertexAlreadyExists {
					logger.Error(err.Error(), zap.String("Edition", gID.Edition))
				}
			}

			_, err = guacIdGraph.Edge("Edition|"+gID.Edition, "Name|"+gID.Name)
			if err != nil {
				err = guacIdGraph.AddEdge("Edition|"+gID.Edition, "Name|"+gID.Name)
				if err != nil && err != graph.ErrEdgeAlreadyExists {
					logger.Error(err.Error())
				}
			}
		}

		if gID.SubPath != "" {

			// add subpath
			_, err := guacIdGraph.Vertex(gID.SubPath)
			if err != nil {
				nodeType := schemas.NodeTypeSoft
				if helpers.IsSHAOrUUID(gID.Name) {
					nodeType = schemas.NodeTypeHard
				}
				err = guacIdGraph.AddVertex(&schemas.GuacIDNode{NodeID: "SubPath|" + gID.SubPath, NodeType: nodeType})
				if err != nil && err != graph.ErrVertexAlreadyExists {
					logger.Error(err.Error(), zap.String("SubPath", gID.SubPath))
				}
			}

			_, err = guacIdGraph.Edge("SubPath|"+gID.SubPath, "Name|"+gID.Name)
			if err != nil {
				err = guacIdGraph.AddEdge("SubPath|"+gID.SubPath, "Name|"+gID.Name)
				if err != nil && err != graph.ErrEdgeAlreadyExists {
					logger.Error(err.Error())
				}
			}
		}

		if gID.Version != "" {
			_, err := guacIdGraph.Vertex(gID.Version)
			if err != nil {
				nodeType := schemas.NodeTypeSoft
				if helpers.IsSHAOrUUID(gID.Name) {
					nodeType = schemas.NodeTypeHard
				}
				err = guacIdGraph.AddVertex(&schemas.GuacIDNode{NodeID: "Version|" + gID.Version, NodeType: nodeType})
				if err != nil && err != graph.ErrVertexAlreadyExists {
					logger.Error(err.Error(), zap.String("Version", gID.Version))
				}
			}

			_, err = guacIdGraph.Edge("Version|"+gID.Version, "Name|"+gID.Name)
			if err != nil {
				err = guacIdGraph.AddEdge("Version|"+gID.Version, "Name|"+gID.Name)
				if err != nil && err != graph.ErrEdgeAlreadyExists {
					logger.Error(err.Error())
				}
			}
		}

		if gID.PkgRel != "" {

			_, err := guacIdGraph.Vertex(gID.PkgRel)
			if err != nil {
				nodeType := schemas.NodeTypeSoft
				if helpers.IsSHAOrUUID(gID.Name) {
					nodeType = schemas.NodeTypeHard
				}
				err = guacIdGraph.AddVertex(&schemas.GuacIDNode{NodeID: "PkgRel|" + gID.PkgRel, NodeType: nodeType})
				if err != nil && err != graph.ErrVertexAlreadyExists {
					logger.Error(err.Error(), zap.String("PkgRel", gID.PkgRel))
				}
			}

			_, err = guacIdGraph.Edge("PkgRel|"+gID.PkgRel, "Name|"+gID.Name)
			if err != nil {
				err = guacIdGraph.AddEdge("PkgRel|"+gID.PkgRel, "Name|"+gID.Name)
				if err != nil && err != graph.ErrEdgeAlreadyExists {
					logger.Error(err.Error())
				}
			}
		}

		if gID.Namespace != "" {

			_, err := guacIdGraph.Vertex(gID.Namespace)
			if err != nil {
				nodeType := schemas.NodeTypeSoft
				if helpers.IsSHAOrUUID(gID.Name) {
					nodeType = schemas.NodeTypeHard
				}
				err = guacIdGraph.AddVertex(&schemas.GuacIDNode{NodeID: "Namespace|" + gID.Namespace, NodeType: nodeType})
				if err != nil && err != graph.ErrVertexAlreadyExists {
					logger.Error(err.Error(), zap.String("Namespace", gID.Namespace))
				}
			}

			_, err = guacIdGraph.Edge("Namespace|"+gID.Namespace, "Name|"+gID.Name)
			if err != nil {
				err = guacIdGraph.AddEdge("Namespace|"+gID.Namespace, "Name|"+gID.Name)
				if err != nil && err != graph.ErrEdgeAlreadyExists {
					logger.Error(err.Error())
				}
			}
		}

		if len(gID.Other) != 0 {
			for _, other := range gID.Other {
				_, err := guacIdGraph.Vertex(other)
				if err != nil {
					nodeType := schemas.NodeTypeSoft
					if helpers.IsSHAOrUUID(gID.Name) {
						nodeType = schemas.NodeTypeHard
					}
					err = guacIdGraph.AddVertex(&schemas.GuacIDNode{NodeID: "Other|" + other, NodeType: nodeType})
					if err != nil && err != graph.ErrVertexAlreadyExists {
						logger.Error(err.Error(), zap.String("Other", other))
					}
				}

				_, err = guacIdGraph.Edge("Other|"+other, "Name|"+gID.Name)
				if err != nil {
					err = guacIdGraph.AddEdge("Other|"+other, "Name|"+gID.Name)
					if err != nil && err != graph.ErrEdgeAlreadyExists {
						logger.Error(err.Error())
					}
				}
			}
		}

	}
	return guacIdGraph, nil
}
