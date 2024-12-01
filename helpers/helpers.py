from typing import Dict, List, Tuple
from schemas import *
import networkx as nx


CPE_components = [
    "part",
    "vendor",
    "product",
    "version",
    "update",
    "edition",
    "language",
    "sw_edition",
    "target_sw",
    "target_hw",
    "other",
]


def parse_cpe(cpe_string: str) -> Dict[str, str]:
    if not cpe_string.startswith("cpe:2.3:"):
        raise ValueError("Invalid CPE string format. Expected CPE 2.3 format.")
    parts = cpe_string.split(":")
    components = parts[2:]

    return components

# In CPE, the product field tells us which software artifact is being referred to.
# In PURL, the name field specifies the software artifact.
def Identifiers_to_nodes_and_edges(
    packages: List['PackageType'], artifacts: List['Artifact'], metadata: List['Metadata']
) -> Tuple[Dict[str, IdentifierNode], List[IdentifierEdge]]:
    dict_of_nodes: Dict[str, IdentifierNode] = {}
    list_of_edges: List[IdentifierEdge] = []
        

    for meta in metadata:
        if meta.key == "cpe":
            components = parse_cpe(cpe_string=meta.value)
            
            if  components[2] != "*" and components[2] not in dict_of_nodes: # since product is the artifact name here
                   dict_of_nodes[components[2]] = (
                        IdentifierNode(id=components[2], type="SymPurl" , attributes={"from_identifier_type": "cpe"})
                    )
                   list_of_edges.append(IdentifierEdge(id="GUACID||"+ components[2], source="GAUCID", target=components[2], attributes={"link_type": "connector"}), )
            
      
            for index, component in enumerate(components):
                
                if index == 2 : 
                   continue # already added to nodes list
                   
                if (
                    component != "*"
                    and CPE_components[index] + "|" + component not in dict_of_nodes
                ):
                    dict_of_nodes[CPE_components[index] + "|" + component] = (
                        IdentifierNode(id=CPE_components[index] + "|" + component, type="label" , attributes={"from_identifier_type": "cpe"})
                    )
                    label_prefix = ""
                    if index >= len(CPE_components):
                        label_prefix = CPE_components[len(CPE_components)- 1]
                    else:
                        label_prefix = CPE_components[index]
                    list_of_edges.append(IdentifierEdge(id=components[2]+"||"+ label_prefix + "|" + component, source=components[2], target=CPE_components[index] + "|" + component, attributes={"link_type": "to_label"}))
                
    return dict_of_nodes, list_of_edges
                                     
                



def Create_identifier_graph(nodes: dict[str, IdentifierNode], edges: List[IdentifierEdge]):

    G = nx.DiGraph()

    # Add nodes to the graph
    for node_id, node in nodes.items():
        G.add_node(node_id, **node.attributes)

    # Add edges to the graph
    for edge in edges:
        G.add_edge(edge.source, edge.target, **edge.attributes)

    return G
