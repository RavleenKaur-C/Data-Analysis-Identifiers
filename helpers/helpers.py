from typing import Dict, List, Tuple
from schemas import *
import networkx as nx

import networkx as nx
from node2vec import Node2Vec
import matplotlib.pyplot as plt
from typing import Dict, List
from scipy.spatial.distance import cosine

import torch
import torch.nn as nn
import torch.optim as optim
import networkx as nx
import numpy as np

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
    packages: List["PackageType"],
    artifacts: List["Artifact"],
    metadata: List["Metadata"],
) -> Tuple[Dict[str, IdentifierNode], List[IdentifierEdge]]:
    dict_of_nodes: Dict[str, IdentifierNode] = {}
    list_of_edges: List[IdentifierEdge] = []

    for meta in metadata:
        if meta.key == "cpe":
            components = parse_cpe(cpe_string=meta.value)

            if (
                components[2] != "*" and components[2] not in dict_of_nodes
            ):  # since product is the artifact name here
                dict_of_nodes[components[2]] = IdentifierNode(
                    id=components[2],
                    type="SymPurl",
                    attributes={"from_identifier_type": "cpe"},
                )
                list_of_edges.append(
                    IdentifierEdge(
                        id="GUACID||" + components[2],
                        source="GAUCID",
                        target=components[2],
                        attributes={"link_type": "connector"},
                    ),
                )

            for index, component in enumerate(components):

                if index == 2:
                    continue  # already added to nodes list

                if (
                    component != "*"
                    and CPE_components[index] + "|" + component not in dict_of_nodes
                ):
                    dict_of_nodes[CPE_components[index] + "|" + component] = (
                        IdentifierNode(
                            id=CPE_components[index] + "|" + component,
                            type="label",
                            attributes={"from_identifier_type": "cpe"},
                        )
                    )
                    label_prefix = ""
                    if index >= len(CPE_components):
                        label_prefix = CPE_components[len(CPE_components) - 1]
                    else:
                        label_prefix = CPE_components[index]
                    list_of_edges.append(
                        IdentifierEdge(
                            id=components[2] + "||" + label_prefix + "|" + component,
                            source=components[2],
                            target=CPE_components[index] + "|" + component,
                            attributes={"link_type": "to_label"},
                        )
                    )

    return dict_of_nodes, list_of_edges


def Create_identifier_graph(
    nodes: dict[str, IdentifierNode], edges: List[IdentifierEdge]
):

    G = nx.DiGraph()

    for node_id, node in nodes.items():
        G.add_node(node_id, **node.attributes)

    for edge in edges:
        G.add_edge(edge.source, edge.target, **edge.attributes)

    return G


def Visualize_graph(G: nx.Graph):
    plt.figure(figsize=(10, 7))
    pos = nx.spring_layout(G, seed=42)

    nx.draw_networkx_nodes(G, pos, node_size=700, node_color="lightblue")

    nx.draw_networkx_edges(G, pos, edge_color="gray")

    nx.draw_networkx_labels(G, pos, font_size=12, font_color="black")

    plt.title("Identifier Graph", fontsize=16)
    plt.show()


def Generate_embeddings_node2vec(
    G: nx.Graph, dimensions: int = 64
) -> Dict[str, List[float]]:
    node2vec = Node2Vec(
        G, dimensions=dimensions, walk_length=30, num_walks=200, workers=1
    )
    model = node2vec.fit(window=10, min_count=1, batch_words=4)

    embeddings = {str(node): model.wv[str(node)].tolist() for node in G.nodes()}
    return embeddings


def create_adjacency_matrix(graph: nx.Graph):
    return nx.to_numpy_array(graph)


class GraphAutoencoder(nn.Module):
    def __init__(self, input_dim, hidden_dim, embedding_dim):
        super(GraphAutoencoder, self).__init__()

        self.encoder = nn.Sequential(
            nn.Linear(input_dim, hidden_dim),
            nn.ReLU(),
            nn.Linear(hidden_dim, embedding_dim),
        )

        self.decoder = nn.Sequential(
            nn.Linear(embedding_dim, hidden_dim),
            nn.ReLU(),
            nn.Linear(hidden_dim, input_dim),
            nn.Sigmoid(),
        )

    def forward(self, adjacency_matrix):

        embeddings = self.encoder(adjacency_matrix)

        reconstructed = self.decoder(embeddings)
        return reconstructed, embeddings


def train_graph_autoencoder(
    graph, hidden_dim=32, embedding_dim=16, epochs=100, lr=0.01
):

    adjacency_matrix = create_adjacency_matrix(graph)
    adjacency_matrix_tensor = torch.tensor(adjacency_matrix, dtype=torch.float32)

    input_dim = adjacency_matrix.shape[1]
    model = GraphAutoencoder(input_dim, hidden_dim, embedding_dim)
    optimizer = optim.Adam(model.parameters(), lr=lr)
    loss_function = nn.MSELoss()

    for epoch in range(epochs):
        optimizer.zero_grad()
        reconstructed, embeddings = model(adjacency_matrix_tensor)
        loss = loss_function(reconstructed, adjacency_matrix_tensor)
        loss.backward()
        optimizer.step()

        if (epoch + 1) % 10 == 0:
            print(f"Epoch [{epoch+1}/{epochs}], Loss: {loss.item():.4f}")

    return model, embeddings


def visualize_embeddings(embeddings, labels=None):
    import matplotlib.pyplot as plt
    from sklearn.decomposition import PCA

    pca = PCA(n_components=2)
    reduced_embeddings = pca.fit_transform(embeddings.detach().numpy())

    plt.figure(figsize=(8, 6))
    plt.scatter(
        reduced_embeddings[:, 0], reduced_embeddings[:, 1], c="blue", label="Nodes"
    )
    if labels:
        for i, label in enumerate(labels):
            plt.annotate(label, (reduced_embeddings[i, 0], reduced_embeddings[i, 1]))
    plt.title(
        "Node Embeddings Visualization (Edge reconstruction-minimize distance based loss)"
    )
    plt.show()
