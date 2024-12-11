from setuptools import setup, find_packages

setup(
    name="dependency-graph-generator",
    version="0.1.0",
    description="A tool to parse SBOM JSON files and generate dependency graphs.",
    author="Ravleen Kaur Chhabra",
    author_email="chhabra4@purdue.edu",
    packages=find_packages(),
    install_requires=[
        "matplotlib",
        "networkx",
    ],
)
