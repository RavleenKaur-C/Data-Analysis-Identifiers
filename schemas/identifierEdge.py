from typing import Dict
import json

class IdentifierEdge:
    def __init__(self, id: str, source: str, target: str, attributes: Dict[str, any]):
        self.id = id
        self.source = source
        self.target = target
        self.attributes = attributes

    def to_json(self) -> str:

        return json.dumps({
            "id": self.id,
            "source": self.source,
            "target": self.target
        }, indent=4) 

    def __repr__(self):

        return f"IdentifierNode(id={self.id}, source={self.source}, target={self.target}, attributes={self.attributes})"
