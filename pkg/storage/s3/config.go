package s3

type s3Config struct {
	accessKey  string
	secretKey  string
	region     string
	bucketName string
}

func NewS3Config(accessKey, secretKey, region, bucketName string) *s3Config {
	return &s3Config{
		accessKey:  accessKey,
		secretKey:  secretKey,
		region:     region,
		bucketName: bucketName,
	}
}
