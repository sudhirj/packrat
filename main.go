package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	port := os.Getenv("PORT")
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("Could not start AWS session")
	}
	s3Svc := s3.New(sess)
	uploader := s3manager.NewUploaderWithClient(s3Svc)

	http.ListenAndServe(":"+port, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			pathParts := strings.Split(r.URL.Path, "/")
			bucket, remainingParts := pathParts[0], pathParts[1:]
			id := r.Header.Get("Logplex-Frame-Id")
			key := strings.Join(append(remainingParts, id), "/")
			uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
				Body:   r.Body,
			})
		}))

}
