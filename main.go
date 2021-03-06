package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	port := os.Getenv("PORT")
	uploader := s3manager.NewUploader(session.New(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)}))

	log.Println("Starting on port:", port)
	http.ListenAndServe(":"+port, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.FormValue("token") != os.Getenv("TOKEN") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
			defer r.Body.Close()
			pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			bucket, remainingParts := pathParts[0], pathParts[1:]
			id := r.Header.Get("Logplex-Frame-Id")
			key := strings.Join(append(remainingParts, id), "/")
			_, err := uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
				Body:   r.Body,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			log.Println(bucket, ":", key)
		}))
}
