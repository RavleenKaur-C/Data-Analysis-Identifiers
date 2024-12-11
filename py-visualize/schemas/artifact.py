import json
from typing import List

class Artifact:
    def __init__(self, id: str, algorithm: str, digest: str):
        self.id = id
        self.algorithm = algorithm
        self.digest = digest

    def __repr__(self):
        return f"Artifact(id={self.id}, algorithm={self.algorithm}, digest={self.digest})"

    @staticmethod
    def from_dict(data: dict) -> 'Artifact':
        return Artifact(
            id=data["id"],
            algorithm=data["algorithm"],
            digest=data["digest"]
        )

class ArtifactParser:
    @staticmethod
    def parse(json_file_path: str) -> List[Artifact]:
        with open(json_file_path, 'r') as file:
            data = json.load(file)
        return [Artifact.from_dict(item) for item in data]

