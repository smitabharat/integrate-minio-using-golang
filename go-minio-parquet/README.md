
# Go Fiber + MinIO Parquet S3 Select

This project demonstrates how to fetch Parquet files from MinIO using S3 Select and return JSON data via a Fiber API.

## Usage

1. Replace MinIO endpoint, accessKey, secretKey, and bucket in `main.go`.
2. Run:

```bash
go mod tidy
go run main.go
```

3. Visit:

```
http://localhost:3000/parquet-data
```

You will get JSON response with data from all Parquet files in the folder.
