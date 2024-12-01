import json
from typing import List, Optional

class PackageVersion:
    def __init__(self, id: str, purl: str, version: str, qualifiers: List[str], subpath: str):
        self.id = id
        self.purl = purl
        self.version = version
        self.qualifiers = qualifiers
        self.subpath = subpath

    def __repr__(self):
        return f"PackageVersion(id={self.id}, purl={self.purl}, version={self.version})"

    @staticmethod
    def from_dict(data: dict) -> 'PackageVersion':
        return PackageVersion(
            id=data["id"],
            purl=data["purl"],
            version=data["version"],
            qualifiers=data["qualifiers"],
            subpath=data["subpath"]
        )

class PackageName:
    def __init__(self, id: str, name: str, versions: List[PackageVersion]):
        self.id = id
        self.name = name
        self.versions = versions

    def __repr__(self):
        return f"PackageName(id={self.id}, name={self.name}, versions={self.versions})"

    @staticmethod
    def from_dict(data: dict) -> 'PackageName':
        return PackageName(
            id=data["id"],
            name=data["name"],
            versions=[PackageVersion.from_dict(v) for v in data["versions"]]
        )

class PackageNamespace:
    def __init__(self, id: str, namespace: str, names: List[PackageName]):
        self.id = id
        self.namespace = namespace
        self.names = names

    def __repr__(self):
        return f"PackageNamespace(id={self.id}, namespace={self.namespace}, names={self.names})"

    @staticmethod
    def from_dict(data: dict) -> 'PackageNamespace':
        return PackageNamespace(
            id=data["id"],
            namespace=data["namespace"],
            names=[PackageName.from_dict(n) for n in data["names"]]
        )

class Subject:
    def __init__(self, id: str, type: str, namespaces: List[PackageNamespace]):
        self.id = id
        self.type = type
        self.namespaces = namespaces

    def __repr__(self):
        return f"Subject(id={self.id}, type={self.type}, namespaces={self.namespaces})"

    @staticmethod
    def from_dict(data: dict) -> 'Subject':
        return Subject(
            id=data["id"],
            type=data["type"],
            namespaces=[PackageNamespace.from_dict(ns) for ns in data["namespaces"]]
        )

class Metadata:
    def __init__(self, id: str, subject: Subject, key: str, value: str, timestamp: str, justification: str, origin: str, collector: str, document_ref: str):
        self.id = id
        self.subject = subject
        self.key = key
        self.value = value
        self.timestamp = timestamp
        self.justification = justification
        self.origin = origin
        self.collector = collector
        self.document_ref = document_ref

    def __repr__(self):
        return f"Metadata(id={self.id}, subject={self.subject}, key={self.key}, value={self.value})"

    @staticmethod
    def from_dict(data: dict) -> 'Metadata':
        return Metadata(
            id=data["id"],
            subject=Subject.from_dict(data["subject"]),
            key=data["key"],
            value=data["value"],
            timestamp=data["timestamp"],
            justification=data["justification"],
            origin=data["origin"],
            collector=data["collector"],
            document_ref=data["documentRef"]
        )

class MetadataParser:
    @staticmethod
    def parse(json_file_path: str) -> List[Metadata]:
        with open(json_file_path, 'r') as file:
            data = json.load(file)
        return [Metadata.from_dict(item) for item in data]
