import psycopg2
from typing import List, Tuple, Any


class PostgresHelper:
    def __init__(self, host: str, database: str, user: str, password: str, port: int = 5432):
        """
        Initialize the PostgresHelper with connection details.
        """
        try:
            self.connection = psycopg2.connect(
                host=host,
                database=database,
                user=user,
                password=password,
                port=port
            )
            self.cursor = self.connection.cursor()
            print("Connected to the PostgreSQL database.")
        except Exception as e:
            print(f"Error connecting to PostgreSQL: {e}")
            raise

    def fetch_edges(self, query: str, params: tuple = None) -> List[Tuple[Any, ...]]:
        """
        Execute a query to fetch data (edges) from the database.
        :param query: The SQL query to execute.
        :param params: Optional tuple of parameters for parameterized queries.
        :return: List of tuples containing the query results.
        """
        try:
            self.cursor.execute(query, params)
            results = self.cursor.fetchall()
            return results
        except Exception as e:
            print(f"Error executing query: {e}")
            return []

    def close_connection(self) -> None:
        """
        Close the database connection.
        """
        try:
            self.cursor.close()
            self.connection.close()
            print("Connection to PostgreSQL closed.")
        except Exception as e:
            print(f"Error closing the connection: {e}")


if __name__ == "__main__":
    #initialize the helper
    postgres_helper = PostgresHelper(
        host="localhost",
        database="your_database",
        user="your_user",
        password="your_password"
    )

    #example query to fetch edges (will replace with the actual query later)
    #the actual query should pull (namespace, name, version, qualifiers)
    query = """
    SELECT namespace, name, version, qualifiers
    FROM your_table_name
    WHERE some_condition = %s
    """
    params = ("your_condition_value",)

    #fetch edges
    edges = postgres_helper.fetch_edges(query, params)
    print("Fetched Edges:", edges)

    #close the connection
    postgres_helper.close_connection()
