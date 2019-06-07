package main

import (
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/s3"
)

// getS3Buckets
func getS3Buckets(sess *session.Session) {
	client := s3.New(sess)
	input := &s3.ListBucketsInput{}
	response, _ := client.ListBuckets(input)
	region := *sess.Config.Region

	var wg sync.WaitGroup

	for _, b := range response.Buckets {
		wg.Add(1)
		go func(c *s3.Bucket) {
			var s3Bucket s3Bucket
			defer wg.Done()
			s3Bucket.name = *c.Name
			s3Bucket.bucketSizeBytes = make(map[string]float64)
			s3Bucket.bucketSizeBytes = map[string]float64{"StandardStorage": 0.0, "StandardIAStorage": 0.0, "ReducedRedundancyStorage": 0.0, "GlacierStorage": 0.0}
			input := &s3.GetBucketLocationInput{Bucket: aws.String(*c.Name)}
			result, _ := client.GetBucketLocation(input)
			if result.LocationConstraint == nil {
				s3Bucket.region = "us-east-1"
			}
			if result.LocationConstraint != nil {
				switch *result.LocationConstraint {
				case "EU":
					s3Bucket.region = "eu-west-1"
				default:
					s3Bucket.region = *result.LocationConstraint
				}
			}

			s3Bucket.checkS3Policy(sess)
			s3Bucket.checkACL(sess)

			if s3Bucket.region == region {
				s3Bucket.getStats(sess)
			} else {
				sessCopy := sess.Copy(&aws.Config{Region: aws.String("us-east-1")})
				s3Bucket.getStats(sessCopy)
			}
			s3Buckets = append(s3Buckets, s3Bucket)
		}(b)
	}
	wg.Wait()
}

func (s3Bucket *s3Bucket) checkS3Policy(sess *session.Session) {
	client := s3.New(sess)
	input := &s3.GetBucketPolicyStatusInput{Bucket: aws.String(s3Bucket.name)}
	response, _ := client.GetBucketPolicyStatus(input)
	if response.PolicyStatus != nil {
		s3Bucket.isPublic = *response.PolicyStatus.IsPublic
	}
}

func (s3Bucket *s3Bucket) checkACL(sess *session.Session) {
	client := s3.New(sess)
	input := &s3.GetBucketAclInput{Bucket: aws.String(s3Bucket.name)}
	response, _ := client.GetBucketAcl(input)
	for _, acl := range response.Grants {
		if (*acl.Permission == "READ" || *acl.Permission == "FULL_CONTROL") && *acl.Grantee.Type == "Group" {
			s3Bucket.isPublic = true
		}
	}
}

// getStats retrieves main statistics of buckets
func (s3Bucket *s3Bucket) getStats(sess *session.Session) {
	client := cloudwatch.New(sess) // New client for CloudWatch
	currentTime := time.Now()
	previousTime := currentTime.AddDate(0, 0, -2)

	queries := []*cloudwatch.MetricDataQuery{}

	query := &cloudwatch.MetricDataQuery{
		Id:    aws.String("m0"),
		Label: aws.String("NumberOfObjects"),
		MetricStat: &cloudwatch.MetricStat{
			Metric: &cloudwatch.Metric{
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("BucketName"),
						Value: aws.String(s3Bucket.name),
					},
					{
						Name:  aws.String("StorageType"),
						Value: aws.String("AllStorageTypes"),
					},
				},
				MetricName: aws.String("NumberOfObjects"),
				Namespace:  aws.String("AWS/S3"),
			},
			Period: aws.Int64(84600),
			Stat:   aws.String("Average"),
		},
	}
	queries = append(queries, query)

	listStorage := []string{"StandardStorage", "StandardIAStorage", "ReducedRedundancyStorage", "GlacierStorage"}
	for s, storageType := range listStorage {
		query := &cloudwatch.MetricDataQuery{
			Id:    aws.String("m" + strconv.Itoa(s+1)),
			Label: aws.String(storageType),
			MetricStat: &cloudwatch.MetricStat{
				Metric: &cloudwatch.Metric{
					Dimensions: []*cloudwatch.Dimension{
						{
							Name:  aws.String("BucketName"),
							Value: aws.String(s3Bucket.name),
						},
						{
							Name:  aws.String("StorageType"),
							Value: aws.String(storageType),
						},
					},
					MetricName: aws.String("BucketSizeBytes"),
					Namespace:  aws.String("AWS/S3"),
				},
				Period: aws.Int64(84600),
				Stat:   aws.String("Average"),
			},
		}
		queries = append(queries, query)
	}

	parameters := &cloudwatch.GetMetricDataInput{
		StartTime:         aws.Time(previousTime),
		EndTime:           aws.Time(currentTime),
		MetricDataQueries: queries,
	}

	var response *cloudwatch.GetMetricDataOutput
	var err error
	var sendReq func()
	sendReq = func() {
		response, err = client.GetMetricData(parameters)
		if err != nil {
			time.Sleep(time.Second * 1)
			sendReq()
		}
	}

	sendReq()

	for _, r := range response.MetricDataResults {
		if *r.Id == "m0" && r.Values != nil {
			s3Bucket.numberOfObjects = *r.Values[0]
		} else if *r.Id != "m0" && r.Values != nil {
			s3Bucket.bucketSizeBytes[*r.Label] = *r.Values[0]
		}
	}
}
