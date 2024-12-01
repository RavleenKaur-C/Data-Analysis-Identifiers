from schemas import *

from helpers import *
import networkx as nx
import os

if __name__ == "__main__":
    packages = PackageParser.parse("data/identifiers/Packages.json")
    metadatas   = MetadataParser.parse("data/identifiers/HasMetadata.json")
    artifacts = ArtifactParser.parse("data/identifiers/Artifacts.json")
    
    nodes, edges = Identifiers_to_nodes_and_edges(packages=packages, metadata=metadatas, artifacts=artifacts)
    
    G = Create_identifier_graph(nodes=nodes, edges=edges)
    # if not nx.is_connected(G): #not implemented for directed type
    #     print("Identifier graph is not connected")
    #     os._exit(1)
    
    Visualize_graph(G=G)
    
        
    
        
        
    
    
   