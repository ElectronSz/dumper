# Dumper

**A powerful and flexible command-line tool for database backups.**

`Dumper` is a Go-based CLI utility that simplifies the process of creating backups for various databases, including PostgreSQL, MySQL, and MongoDB. It's built on top of the `dump_util` library and offers features like compression, batch processing, and table exclusion to ensure efficient and customized database dumps.

## Features

-   **Multi-Database Support**: Works with PostgreSQL, MySQL, and MongoDB.
-   **Efficient Dumping**: Uses concurrent workers and batch processing for fast backups.
-   **Output Compression**: Supports Gzip compression to save disk space.
-   **Table Exclusion**: Easily exclude specific tables or collections from your dump.
-   **Customizable**: Adjust batch size and the number of workers to fit your needs.
