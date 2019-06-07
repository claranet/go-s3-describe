package main

import (
	"flag"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gosuri/uilive"
)

var s3Buckets []s3Bucket
var wg sync.WaitGroup

func init() {
	s3Buckets = make([]s3Bucket, 0)
}

func main() {
	profile := flag.String("profile", "default", "Profile from ~/.aws/config")
	region := flag.String("region", "eu-west-1", "Region  (only to create session)")
	flag.Parse()

	// Create session (credentials from ~/.aws/config)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState:       session.SharedConfigEnable,  //enable use of ~/.aws/config
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider, //ask for MFA if needed
		Profile:                 string(*profile),
		Config:                  aws.Config{Region: aws.String(*region)},
	}))

	// Progress Bar
	progressBar := uilive.New()
	progressBar.Start()
	go func(progressBar *uilive.Writer) {
		for {
			fmt.Fprintf(progressBar, "In progress\n")
			time.Sleep(time.Millisecond * 500)
			fmt.Fprintf(progressBar, "In progress.\n")
			time.Sleep(time.Millisecond * 500)
			fmt.Fprintf(progressBar, "In progress..\n")
			time.Sleep(time.Millisecond * 500)
			fmt.Fprintf(progressBar, "In progress...\n")
			time.Sleep(time.Millisecond * 500)
			fmt.Fprintf(progressBar, "              \n")
			time.Sleep(time.Millisecond * 5)
		}
	}(progressBar)

	getS3Buckets(sess)
	wg.Wait()

	fmt.Fprintf(progressBar, "                   \n\n")
	progressBar.Stop()

	if len(s3Buckets) != 0 {
		sort.Sort(bySize(s3Buckets))
		PrintResult(&s3Buckets)
	} else {
		fmt.Println("No bucket found")
	}
}
