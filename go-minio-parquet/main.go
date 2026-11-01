package main

import (
	"bytes"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/xitongsys/parquet-go-source/buffer"
	"github.com/xitongsys/parquet-go/reader"
)

type Record struct {
	ConsumerNo     string `parquet:"name=Consumer_No, type=BYTE_ARRAY, convertedtype=UTF8"`
	BillNo         string `parquet:"name=Bill_No, type=BYTE_ARRAY, convertedtype=UTF8"`
	Name           string `parquet:"name=Name, type=BYTE_ARRAY, convertedtype=UTF8"`
	Address        string `parquet:"name=Address, type=BYTE_ARRAY, convertedtype=UTF8"`
	BillDate       string `parquet:"name=Bill_Date, type=BYTE_ARRAY, convertedtype=UTF8"`
	DueDate        string `parquet:"name=Due_Date, type=BYTE_ARRAY, convertedtype=UTF8"`
	UnitsConsumed  string `parquet:"name=Units_Consumed, type=BYTE_ARRAY, convertedtype=UTF8"`
	RatePerUnit    string `parquet:"name=Rate_per_Unit, type=BYTE_ARRAY, convertedtype=UTF8"`
	TotalAmountDue string `parquet:"name=Total_Amount_Due, type=BYTE_ARRAY, convertedtype=UTF8"`
	File           string `parquet:"name=File, type=BYTE_ARRAY, convertedtype=UTF8"`
	ProcessedAt    string `parquet:"name=Processed_At, type=BYTE_ARRAY, convertedtype=UTF8"`
}

func main() {
	app := fiber.New()
	app.Get("/parquet-data", func(c *fiber.Ctx) error {
		endpoint := "13.126.123.126:9000"
		accessKey := "admin"
		secretKey := "admin12345"
		bucket := "remp-test-bucket"
		prefix := "ocr_meter_reading/"

		s3Client, _ := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: false,
		})

		ctx := context.Background()
		objectCh := s3Client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true})

		allData := map[string][]Record{}

		for obj := range objectCh {
			if obj.Err != nil {
				continue
			}

			object, _ := s3Client.GetObject(ctx, bucket, obj.Key, minio.GetObjectOptions{})
			buf := new(bytes.Buffer)
			buf.ReadFrom(object)
			object.Close()

			fr := buffer.NewBufferFile()
			fr.Write(buf.Bytes())
			pr, err := reader.NewParquetReader(fr, new(Record), 4)
			if err != nil {
				log.Println("Failed to create reader:", err)
				continue
			}
			num := int(pr.GetNumRows())
			records := make([]Record, num)
			if err := pr.Read(&records); err != nil {
				log.Println("Failed to read parquet:", err)
				pr.ReadStop()
				continue
			}
			pr.ReadStop()

			// Example: filter data like a SQL WHERE
			var filtered []Record
			for _, r := range records {
				if r.ConsumerNo == "46021000001" { // filter condition
					filtered = append(filtered, r)
				}
			}

			allData[obj.Key] = filtered
		}

		return c.JSON(fiber.Map{"files_data": allData})
	})

	log.Fatal(app.Listen(":3000"))
}
