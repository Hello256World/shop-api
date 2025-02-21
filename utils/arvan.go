package utils

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	maxPartSize        = int64(400 * 1024 * 1024)
	maxRetries         = 3
	awsAccessKeyID     = "3c5421dc-a54f-4788-9164-6816a2e2d323"
	awsSecretAccessKey = "3d04700e754d2376902e02f27b8f2b6e78eef33c4f0cbca971feb416ec682133"
	awsBucketRegion    = "default"
	awsBucketEndpoint  = "s3.ir-thr-at1.arvanstorage.ir"
)

func uploadToArvan(bucketName string, folder string, imagePath string) *string {
	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")
	_, err := creds.Get()
	if err != nil {
		fmt.Printf("bad credentials: %s", err)
	}
	cfg := aws.NewConfig().WithRegion(awsBucketRegion).WithCredentials(creds).WithEndpoint(awsBucketEndpoint)
	session, _ := session.NewSession()
	svc := s3.New(session, cfg)

	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Printf("err opening file: %s", err)
		return nil
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size)
	fileType := http.DetectContentType(buffer)
	file.Read(buffer)

	newFileName := uuid.New().String() + filepath.Ext(file.Name())
	path := "/" + folder + "/" + newFileName
	input := &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(path),
		ContentType: aws.String(fileType),
		ACL:         aws.String("public-read"),
	}

	resp, err := svc.CreateMultipartUpload(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fmt.Println("Created multipart upload request")

	var curr, partLength int64
	var remaining = size
	var completedParts []*s3.CompletedPart
	partNumber := 1
	for curr = 0; remaining != 0; curr += partLength {
		if remaining < maxPartSize {
			partLength = remaining
		} else {
			partLength = maxPartSize
		}
		completedPart, err := uploadPart(svc, resp, buffer[curr:curr+partLength], partNumber)
		if err != nil {
			fmt.Println(err.Error())
			err := abortMultipartUpload(svc, resp)
			if err != nil {
				fmt.Println(err.Error())
			}
			return nil
		}
		remaining -= partLength
		partNumber++
		completedParts = append(completedParts, completedPart)
	}

	completeResponse, err := completeMultipartUpload(svc, resp, completedParts)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	fmt.Printf("Successfully uploaded file: %s\n", completeResponse.String())
	imageLoc := fmt.Sprintf("https://%v", *completeResponse.Location)
	return &imageLoc
}

func completeMultipartUpload(svc *s3.S3, resp *s3.CreateMultipartUploadOutput, completedParts []*s3.CompletedPart) (*s3.CompleteMultipartUploadOutput, error) {
	completeInput := &s3.CompleteMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}
	return svc.CompleteMultipartUpload(completeInput)
}

func uploadPart(svc *s3.S3, resp *s3.CreateMultipartUploadOutput, fileBytes []byte, partNumber int) (*s3.CompletedPart, error) {
	tryNum := 1
	partInput := &s3.UploadPartInput{
		Body:          bytes.NewReader(fileBytes),
		Bucket:        resp.Bucket,
		Key:           resp.Key,
		PartNumber:    aws.Int64(int64(partNumber)),
		UploadId:      resp.UploadId,
		ContentLength: aws.Int64(int64(len(fileBytes))),
	}

	for tryNum <= maxRetries {
		uploadResult, err := svc.UploadPart(partInput)
		if err != nil {
			if tryNum == maxRetries {
				if aerr, ok := err.(awserr.Error); ok {
					return nil, aerr
				}
				return nil, err
			}
			fmt.Printf("Retrying to upload part #%v\n", partNumber)
			tryNum++
		} else {
			fmt.Printf("Uploaded part #%v\n", partNumber)
			return &s3.CompletedPart{
				ETag:       uploadResult.ETag,
				PartNumber: aws.Int64(int64(partNumber)),
			}, nil
		}
	}
	return nil, nil
}

func abortMultipartUpload(svc *s3.S3, resp *s3.CreateMultipartUploadOutput) error {
	fmt.Println("Aborting multipart upload for UploadId#" + *resp.UploadId)
	abortInput := &s3.AbortMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
	}
	_, err := svc.AbortMultipartUpload(abortInput)
	return err
}

func DeleteImageOfServer(bucket, image string) bool {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, _ := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})
	svc := s3.New(sess, &aws.Config{
		Region:   aws.String("default"),
		Endpoint: aws.String(awsBucketEndpoint),
	})

	// Delete the item
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(image)})
	if err != nil {
		exitErrorf("Unable to delete object %q from bucket %q, %v", image, bucket, err)
		return false
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(image),
	})
	if err != nil {
		exitErrorf("Error occurred while waiting for object %q to be deleted, %v", image, err)
		return false
	}

	fmt.Printf("Object %q successfully deleted\n", image)
	return true
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func AddImageToServer(c *gin.Context, bucketName, folder string, file *multipart.FileHeader) (*string, error) {

	fileExt := filepath.Ext(file.Filename)
	if fileExt == "" {
		fileExt = ".jpg"
	}

	destPath := filepath.Join("uploads", "upload-"+fileExt)

	if err := c.SaveUploadedFile(file, destPath); err != nil {
		return nil, errors.New("مشکلی در ذخیره عکس پیش آمده")
	}

	imageLocation := uploadToArvan(bucketName, folder, destPath)

	if imageLocation == nil {
		return nil, errors.New("مشکلی در ذخیره عکس در سرور پیش آمده")

	}

	if err := os.RemoveAll("uploads"); err != nil {
		return nil, err
	}

	return imageLocation, nil
}
