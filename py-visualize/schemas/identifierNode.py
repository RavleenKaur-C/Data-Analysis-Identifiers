
from typing import Dict

import json

class IdentifierNode:
    def __init__(self, id: str, type: str, attributes: dict[str, any]):
        self.id = id
        self.type = type
        self.attributes = attributes

    def to_json(self) -> str:

        return json.dumps({
            "id": self.id,
            "type": self.type,
            "attributes": self.attributes
        }, indent=4) 

    def __repr__(self):
        return f"IdentifierNode(id={self.id}, type={self.type}, attributes={self.attributes})"