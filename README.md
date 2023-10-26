# CSV Processor and Database Loader

## Overview

This repository contains a small script designed to process plain text transaction records and load them into a database. The script is written in the Go programming language (Golang) and provides a streamlined solution for cleaning and importing CSV data into your database.

## Features

- Efficiently processes CSV files.
- Cleans and prepares the data for database insertion.
- Simplifies the data loading process into a database.
- Written in Go for performance and reliability.

## Getting Started

To get started with this CSV processing and database loading script, follow these steps:

1. **Clone the Repository**: 
   ```shell
   git clone https://github.com/jstoledano/CECyRDIngestionTool.git
   ```

2. **Run the Script**:
   ```shell
   go run main.go
   ```

3. **Configure Database Connection**:
   Make sure to update the script with your database connection details in the configuration section.

4. **Provide CSV Files**:
   Place your CSV files in the designated directory.

5. **Execute the Script**:
   Run the script to process and load the CSV data into your database.

## Configuration

Before running the script, configure the database connection details in the `config.go` file.

```go
const (
    DBUser     = "your-username"
    DBPassword = "your-password"
    DBName     = "your-database-name"
)
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

This project was developed to simplify the process of cleaning and loading CSV data into a database using the Go programming language.

Happy Coding!
