package schemas

import (
	"fmt"
	"strings"
)

func ConvertPurlToGuacID(purl Purl) GuacID {
	id := GuacID{
		Ecosystem: purl.Type,
		Namespace: purl.Namespace,
		Name:      purl.Name,
		Version:   purl.Version,

		SubPath: purl.SubPath,
	}

	if purl.Type == "deb" {
		id.PkgRel = purl.Version
	}

	if purl.Type == "core" {
		id.Other = []string{purl.QualX}
	}

	if purl.Type != "cargo" {
		id.Arch = purl.QualArch
	}

	return id
}

func ConvertCPEToGuacID(cpe CPE) GuacID {
	return GuacID{
		Ecosystem: cpe.TargetSW,
		Namespace: cpe.Vendor,
		Name:      cpe.Product,
		Version:   cpe.Version,
		Arch:      cpe.TargetHW,
		Other:     cpe.Other,
		PkgRel:    cpe.Update,
		Edition:   cpe.Edition,
		//subpath does not exist for CPE
	}
}

func ParseCPE(cpeStr string) (CPE, error) {
	// Format: cpe:/<part>:<vendor>:<product>:<version>:<update>:<edition>:<language>:<sw_edition>:<other>:<other>...
	parts := strings.Split(cpeStr, ":")
	if len(parts) < 6 || parts[1] == "*" || parts[2] == "*" || parts[3] == "*" || parts[4] == "*" || parts[5] == "*" {
		return CPE{}, fmt.Errorf("invalid CPE format: %s", cpeStr)
	}

	targetSW := parts[1]
	vendor := parts[2]
	product := parts[3]
	version := parts[4]
	update := parts[5]

	edition := ""
	language := ""
	swEdition := ""

	other := []string{}

	if len(parts) > 6 && parts[6] != "*" {
		edition = parts[6]
	}

	if len(parts) > 7 && parts[7] != "*" {
		language = parts[7]
	}

	if len(parts) > 8 && parts[8] != "*" {
		swEdition = parts[8]
	}

	if len(parts) > 9 {
		for _, val := range parts[9:] {
			if val == "*" {
				continue
			}
			other = append(other, val)
		}
	}

	return CPE{
		TargetSW:  targetSW,
		Vendor:    vendor,
		Product:   product,
		Version:   version,
		Update:    update,
		Edition:   edition,
		Language:  language,
		SWEdition: swEdition,
		Other:     other,
	}, nil
}

func GuacIDNodeID(node *GuacIDNode) string {
	return node.NodeID
}
