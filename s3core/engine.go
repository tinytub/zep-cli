package s3core

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func NewClient(endpoint string, key string, secret string) *s3.S3 {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("default"),
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String(endpoint),
	})
	if err != nil {
		fmt.Printf("Cannot create new s3 seesion, %v", err)
	}
	client := s3.New(sess)
	return client
}

//create bucket
func CreateBucket(svc *s3.S3, key string) {
	params := &s3.CreateBucketInput{
		Bucket: aws.String(key), // Required
		//ACL:    aws.String("BucketCannedACL"),
		//CreateBucketConfiguration: &s3.CreateBucketConfiguration{
		//LocationConstraint: aws.String("sh-bt-1"),
		//},
		//GrantFullControl: aws.String("GrantFullControl"),
		//GrantRead:        aws.String("GrantRead"),
		//GrantReadACP:     aws.String("GrantReadACP"),
		//GrantWrite:       aws.String("GrantWrite"),
		//GrantWriteACP:    aws.String("GrantWriteACP"),
	}

	resp, err := svc.CreateBucket(params)

	if err != nil {
		fmt.Println("create err: ", err.Error())
		return
	}
	fmt.Printf("create done! resp: %s", resp)
	return
}

func ListBucket(svc *s3.S3) {
	result, err := svc.ListBuckets(nil)
	if err != nil {
		fmt.Printf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets: ", len(result.Buckets))
	fmt.Println("Buckets:")
	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

func SetOBJ(svc *s3.S3, bucket, key, value, filename string) (bool, error) {
	//bucketname := bucket

	var putinput *s3.PutObjectInput
	if filename != "" {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Println(err)
			fmt.Printf("%s dosen't exist", filename)
			os.Exit(0)
		}
		value = filename
		file, _ := os.Open(filename)
		defer file.Close()
		putinput = &s3.PutObjectInput{
			//Body:   strings.NewReader("Hi S3"),
			Body:   file,
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}
	} else if value != "" {
		putinput = &s3.PutObjectInput{
			//Body:   strings.NewReader("Hi S3"),
			Body:   strings.NewReader(value),
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}
	}

	_, errput := svc.PutObject(putinput)

	if errput != nil {
		//log.Printf("Failed to upload data to %s/%s, %s\n", bucketname, key, errput)
		return false, errput
	}

	//fmt.Println("put obj result result: ", putobjResult)
	//fmt.Printf("Successfully upload %q to %q/%s\n", value, bucketname, key)
	return true, nil
}

func GetOBJ(svc *s3.S3, bucket, key, output string) (string, error) {
	getobj := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	//resultget, err := svc.GetObjectRequest(getobj)
	resultget, err := svc.GetObject(getobj)
	if err != nil {
		return "", err
	}
	if output != "stdout" {
		file, err := os.Create(output)
		defer file.Close()
		if err != nil {
			//log.Fatal("Failed to create file", err)
			return "", err
		}
		if _, err := io.Copy(file, resultget.Body); err != nil {
			resultget.Body.Close()
			//log.Fatal("Failed to copy object to file", err)
			return "", err
		}
	}
	//io.Copy(os.Stdout, resultget.Body)
	bodyBytes, err := ioutil.ReadAll(resultget.Body)
	resultget.Body.Close()
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)
	return bodyString, nil
	//fmt.Printf("\nThe stdout of key: %s ---> %s \n", key, bodyString)
	/* 这里可以把 body 拷贝回去,用 io.Copy 重新读取
	resultget.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	io.Copy(os.Stdout, resultget.Body)
	*/

	//io.Copy(os.Stdout, resultget.HTTPResponse.Body)
}

func ListOBJ(svc *s3.S3, bucket string) {
	listobj := &s3.ListObjectsInput{
		Bucket: &bucket,
	}
	i := 0
	err := svc.ListObjectsPages(listobj, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		fmt.Println("Page,", i)
		i++
		for _, obj := range p.Contents {
			fmt.Println("Object:", *obj.Key)
		}
		return true
	})
	if err != nil {
		fmt.Println("failed to list objects", err)
	}
}

func DelOBJ(svc *s3.S3, bucket, key string) (bool, error) {
	delobj := &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	_, err := svc.DeleteObject(delobj)
	if err != nil {
		return false, err
	}
	return true, nil
}
