package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	sess := session.Must(session.NewSession()) // The session for S3 to use
	uploader := s3manager.NewUploader(sess)    // Create an uploader with default uploader with session

	picturePath := "Froogo.png"              // Define path to file
	pictureFile, err := os.Open(picturePath) // Read file to bytes
	if err != nil {
		log.Printf("Failed to open file %v: %v.\n", picturePath, err)
	}

	// Upload file to S3
	result, uplErr := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("s.froogo.co.uk"),                 // Bucket name to upload (not necessarily domain)
		Key:    aws.String("Directory/To/Upload/Froogo.png"), // Directory to upload in S3
		Body:   pictureFile,                                  // Body to upload (just bytes)
		// ACL:    aws.String("public-read"),                 // Set to public read (no key required to read)
	})
	if uplErr != nil {
		log.Printf("Failed to upload file: %v.\n", uplErr)
		return
	}
	log.Printf("Picture uploaded to %v!\n", result.Location) // Print URL location to console

	// Generate pre-signed URL
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("s.froogo.co.uk"),
		Key:    aws.String("Directory/To/Upload/Froogo.png"),
	})

	signedURL, sigErr := req.Presign(5 * time.Minute)
	if sigErr != nil {
		log.Printf("Failed to presign object: %v.\n", sigErr)
		return
	}

	fmt.Printf("Presigned URL for picture: %v!\n", signedURL)
}
