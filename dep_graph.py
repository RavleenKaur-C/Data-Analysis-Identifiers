import networkx as nx
import matplotlib.pyplot as plt
import json

#json file
with open("hasSBOM-syft-spdx-k8s.gcr.io-kube-apiserver.v1.24.1.json", "r") as f:
    data = json.load(f)

#directed graph -> dep graph
G = nx.DiGraph()

#sets for nodes and edges
software_nodes = set()
dependency_nodes = set()
edges = []

#processing
for entry in data:
    #software node
    if "subject" in entry and "namespaces" in entry["subject"]:
        namespaces = entry["subject"]["namespaces"]
        if namespaces and "names" in namespaces[0] and namespaces[0]["names"]:
            software_name = namespaces[0]["names"][0]["name"]
            software_nodes.add(software_name)
            
            #dependencies
            if "includedSoftware" in entry:
                for dep in entry["includedSoftware"]:
                    if "namespaces" in dep and dep["namespaces"]:
                        dep_namespace = dep["namespaces"][0]
                        if "names" in dep_namespace and dep_namespace["names"]:
                            dependency_name = dep_namespace["names"][0]["name"]
                            dependency_nodes.add(dependency_name)
                            #edge from software to dep node
                            edges.append((software_name, dependency_name))

dependency_nodes = dependency_nodes - software_nodes

#building graph
G.add_nodes_from(software_nodes, type="software")
G.add_nodes_from(dependency_nodes, type="dependency")
G.add_edges_from(edges)

#removing self loops
G.remove_edges_from(nx.selfloop_edges(G))

#analzing
#coreness
coreness = nx.core_number(G.to_undirected())

#centrality measures
degree_centrality = nx.degree_centrality(G)

#corness and centrality
'''
print("Coreness:")
for node, core in coreness.items():
    print(f"{node}: {core}")

'''
print("\nDegree Centrality:")
for node, centrality in degree_centrality.items():
    print(f"{node}: {centrality:.4f}")


# Visualize the graph
plt.figure(figsize=(12, 8))
pos = nx.spring_layout(G, seed=42) #spring layout
nx.draw_networkx_nodes(G, pos, nodelist=list(software_nodes), node_color="blue", label="Software", node_size=500)
nx.draw_networkx_nodes(G, pos, nodelist=list(dependency_nodes), node_color="orange", label="Dependency", node_size=300)
nx.draw_networkx_edges(G, pos, edgelist=edges, alpha=0.7)
nx.draw_networkx_labels(G, pos, font_size=8, font_color="black")

plt.title("Software Dependency Graph", fontsize=14)
plt.legend(scatterpoints=1)
plt.show()

'''

print("Total nodes:", len(G.nodes))
print("Software nodes:", len(software_nodes))
print("Dependency nodes:", len(dependency_nodes))


#checking to ensure softwares do depend on each other
software_nodes = ["kube-apiserver-v1.24.2", "kube-apiserver-v1.24.1"]

#check dependencies
is_node1_depends_on_node2 = G.has_edge(software_nodes[0], software_nodes[1])
is_node2_depends_on_node1 = G.has_edge(software_nodes[1], software_nodes[0])

print(f"{software_nodes[0]} depends on {software_nodes[1]}: {is_node1_depends_on_node2}")
print(f"{software_nodes[1]} depends on {software_nodes[0]}: {is_node2_depends_on_node1}")


print("Software nodes in dependency_nodes:", set(software_nodes) & dependency_nodes)
'''

