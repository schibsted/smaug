package credentials

type SmaugCredentials struct {
	RoleArn         string `json:"RoleArn"`
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"Token"`
	Expiration      string `json:"Expiration"`
}
