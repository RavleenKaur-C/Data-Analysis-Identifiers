package processidentifiers

import (
	"fmt"
	"math"

	"go-query/schemas"

	"github.com/dominikbraun/graph"
	"gonum.org/v1/gonum/mat"
)

func ComputeModularityMatrix(g graph.Graph[string, *schemas.GuacIDNode]) (*mat.Dense, []string) {
	adjMap, _ := g.AdjacencyMap()
	nodes := make([]string, 0, len(adjMap))
	for nodeID := range adjMap {
		nodes = append(nodes, nodeID)
	}

	
	inDegrees := make(map[string]float64)
	outDegrees := make(map[string]float64)
	totalEdges := 0.0

	for u, neighbors := range adjMap {
		outDegrees[u] = float64(len(neighbors))
		for v := range neighbors {
			inDegrees[v]++
			totalEdges++
		}
	}

	
	n := len(nodes)
	B := mat.NewDense(n, n, nil)

	for i, u := range nodes {
		for j, v := range nodes {
			Auv := 0.0
			if _, exists := adjMap[u][v]; exists {
				Auv = 1.0 
			}
			expected := (outDegrees[u] * inDegrees[v]) / totalEdges
			B.Set(i, j, Auv-expected)
		}
	}

	return B, nodes
}

func SpectralDivision(B *mat.Dense, nodes []string) ([]string, []string, bool) {
	var eigen mat.Eigen
	if ok := eigen.Factorize(B, mat.EigenRight); !ok {
		return nil, nil, false
	}


	eigenvalues := eigen.Values(nil)
	maxIdx, maxVal := -1, math.Inf(-1)
	for i, val := range eigenvalues {
		if real(val) > maxVal {
			maxIdx, maxVal = i, real(val)
		}
	}


	if maxVal < 1e-6 {
		return nil, nil, false
	}


	vecs := mat.NewCDense(len(nodes), len(nodes), nil)
	eigen.VectorsTo(vecs)

	group1, group2 := []string{}, []string{}
	for i, node := range nodes {
		if real(vecs.At(i, maxIdx)) > 0 {
			group1 = append(group1, node)
		} else {
			group2 = append(group2, node)
		}
	}
	return group1, group2, true
}

func RecursiveCommunityDetection(g graph.Graph[string, *schemas.GuacIDNode]) []schemas.Community {
	var communities []schemas.Community
	communityCounter := 1

	var detect func(subgraph graph.Graph[string, *schemas.GuacIDNode], nodes []string, communityID string)
	detect = func(subgraph graph.Graph[string, *schemas.GuacIDNode], nodes []string, communityID string) {
		if len(nodes) <= 1 {
		
			communities = append(communities, schemas.Community{
				CommunityID: communityID,
			
				Size:        len(nodes),
				GraphSubset: &subgraph,
			})
			return
		}


		B, _ := ComputeModularityMatrix(subgraph)

	
		group1, group2, split := SpectralDivision(B, nodes)
		if !split {
			communities = append(communities, schemas.Community{
				CommunityID: communityID,
	
				Size:        len(nodes),
				GraphSubset: &subgraph,
			})
			return
		}

		
		subgraph1 := createSubgraph(subgraph, group1)
		subgraph2 := createSubgraph(subgraph, group2)
		detect(subgraph1, group1, fmt.Sprintf("%s-1", communityID))
		detect(subgraph2, group2, fmt.Sprintf("%s-2", communityID))
	}


	nodes := []string{}
	adjMap, _ := g.AdjacencyMap()
	for nodeID := range adjMap {
		nodes = append(nodes, nodeID)
	}

	detect(g, nodes, fmt.Sprintf("C%d", communityCounter))
	return communities
}


func createSubgraph(g graph.Graph[string, *schemas.GuacIDNode], nodes []string) graph.Graph[string, *schemas.GuacIDNode] {
	subgraph := graph.New[string, *schemas.GuacIDNode](func(n *schemas.GuacIDNode) string { return n.NodeID }, graph.Directed())
	for _, nodeID := range nodes {
		node, _ := g.Vertex(nodeID)
		_ = subgraph.AddVertex(node)
	}

	adjMap, _ := g.AdjacencyMap()
	for u, neighbors := range adjMap {
		if contains(nodes, u) {
			for v := range neighbors {
				if contains(nodes, v) {
					_ = subgraph.AddEdge(u, v)
				}
			}
		}
	}
	return subgraph
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
