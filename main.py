from schemas import *

from helpers import *
import networkx as nx
import os

if __name__ == "__main__":
    packages = PackageParser.parse("data/identifiers/Packages.json")
    metadatas   = MetadataParser.parse("data/identifiers/HasMetadata.json")
    artifacts = ArtifactParser.parse("data/identifiers/Artifacts.json")
    
    nodes_frag, edges_frag = Identifiers_to_nodes_and_edges(packages=packages, metadata=metadatas, artifacts=artifacts, fragmented=True, only_CPE=True, only_Purl=True)
    nodes_conn, edges_conn = Identifiers_to_nodes_and_edges(packages=packages, metadata=metadatas, artifacts=artifacts, fragmented=False, only_CPE=True, only_Purl=True)
    
    G_frag = Create_identifier_graph(nodes=nodes_frag, edges=edges_frag)
    # G_conn = Create_identifier_graph(nodes=nodes_conn, edges=edges_conn)

    
    # Visualize_graph(G=G_conn)
    Visualize_graph(G=G_frag)
    in_degrees = dict(G_frag.in_degree())

    top_10_nodes = sorted(in_degrees.items(), key=lambda x: x[1], reverse=True)

    print("Top 10 nodes by in-degree:")
    for node, degree in top_10_nodes:
        print(f"Node {node}: In-degree {degree}")
        
        
        
    
        
        
    
    
   