import json
from typing import List

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

class PackageType:
    def __init__(self, id: str, type: str, namespaces: List[PackageNamespace]):
        self.id = id
        self.type = type
        self.namespaces = namespaces

    def __repr__(self):
        return f"PackageType(id={self.id}, type={self.type}, namespaces={self.namespaces})"

    @staticmethod
    def from_dict(data: dict) -> 'PackageType':
        return PackageType(
            id=data["id"],
            type=data["type"],
            namespaces=[PackageNamespace.from_dict(ns) for ns in data["namespaces"]]
        )

class PackageParser:
    @staticmethod
    def parse(json_file_path: str) -> List[PackageType]:
        with open(json_file_path, 'r') as file:
            data = json.load(file)
        return [PackageType.from_dict(item) for item in data]
