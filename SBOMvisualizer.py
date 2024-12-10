import json
import os
import networkx as nx
import matplotlib.pyplot as plt

#sbom json file 
def process_sbom_file(file_path, graph, software_nodes, dependency_nodes, edges, software_mapping):
    with open(file_path, "r") as f:
        data = json.load(f)

    #if includedDependencies is empty, skip
    if not data.get("includedDependencies"):
        print(f"Skipping {file_path} (no includedDependencies found)...")
        return

    file_software_nodes = set()  #tracking software nodes

    #extracting software node
    subject = data.get("subject", {})
    if "namespaces" in subject:
        for namespace in subject["namespaces"]:
            if "names" in namespace:
                for name_entry in namespace["names"]:
                    software_name = name_entry.get("name", "Unknown")
                    software_nodes.add(software_name)
                    file_software_nodes.add(software_name)

    #extracting dep nodes
    for dependency_entry in data["includedDependencies"]:
        dependency_package = dependency_entry.get("dependencyPackage", {})
        if "namespaces" in dependency_package:
            for namespace in dependency_package["namespaces"]:
                if "names" in namespace:
                    for name_entry in namespace["names"]:
                        dependency_name = name_entry.get("name", "Unknown")
                        dependency_nodes.add(dependency_name)

                        #adding edge
                        for software_name in file_software_nodes:
                            edges.append((software_name, dependency_name))

    #for logging
    software_mapping[file_path] = file_software_nodes

def main(directory):
    G = nx.DiGraph()

    software_nodes = set()
    dependency_nodes = set()
    edges = []

    software_mapping = {}

    for filename in os.listdir(directory):
        if filename.endswith(".json"):
            file_path = os.path.join(directory, filename)
            print(f"Processing {file_path}...")
            process_sbom_file(file_path, G, software_nodes, dependency_nodes, edges, software_mapping)

    #building graph
    G.add_nodes_from(software_nodes, type="software")
    G.add_nodes_from(dependency_nodes, type="dependency")
    G.add_edges_from(edges)

    G.remove_edges_from(nx.selfloop_edges(G))

    #coreness
    coreness = nx.core_number(G.to_undirected())

    #centrality
    degree_centrality = nx.degree_centrality(G)

    #nodes
    print(f"\nTotal Nodes: {len(G.nodes)}")
    print(f"Software Nodes: {len(software_nodes)}")
    print(f"Dependency Nodes: {len(dependency_nodes)}")
    print(f"Total Edges: {len(G.edges)}")

    #software node names for verification
    print("\nSoftware Names by SBOM File:")
    for file, software in software_mapping.items():
        print(f"{file}: {', '.join(software)}")

    print("\nCentrality and Coreness for Dependency Nodes:")
    for node in dependency_nodes:
        print(f"Dependency Node: {node}")
        print(f"  Degree Centrality: {degree_centrality.get(node, 0):.4f}")
        print(f"  Coreness: {coreness.get(node, 0)}")

    #visualization
    plt.figure(figsize=(12, 8))
    pos = nx.spring_layout(G, seed=42)  #spring layout
    nx.draw_networkx_nodes(G, pos, nodelist=list(software_nodes), node_color="blue", label="Software", node_size=500)
    nx.draw_networkx_nodes(G, pos, nodelist=list(dependency_nodes), node_color="orange", label="Dependency", node_size=300)
    nx.draw_networkx_edges(G, pos, edgelist=edges, alpha=0.7)
    nx.draw_networkx_labels(G, pos, font_size=8, font_color="black")

    plt.title("Software Dependency Graph", fontsize=14)
    plt.legend(scatterpoints=1)
    plt.show()

directory_path = "./data/guac-data-hassboms"  
main(directory_path)
