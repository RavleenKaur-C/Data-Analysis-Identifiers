package processidentifiers

import (
	"context"
	"fmt"

	"github.com/guacsec/guac/pkg/assembler/backends"
	entbackend "github.com/guacsec/guac/pkg/assembler/backends/ent/backend"
	"github.com/guacsec/guac/pkg/assembler/graphql/model"
)

func SetupEntBackendForIdentifiers() (backends.Backend, context.Context, error) {

	ctx := context.Background()

	client, err := entbackend.SetupBackend(ctx, &entbackend.BackendOptions{
		DriverName:  "postgres",
		Address:     "postgres://guac:guac@localhost/guac?sslmode=disable",
		Debug:       false,
		AutoMigrate: true,
	})

	be, err := entbackend.GetBackend(client)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get ent backend %s", err)
	}
	return be, ctx, nil
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
