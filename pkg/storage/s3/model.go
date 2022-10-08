package s3

import "time"

type Media struct {
	Filename     string
	Date         string
	Url          string
	LastModified time.Time
}
