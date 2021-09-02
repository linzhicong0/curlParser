package curlParser

// The struct is to describe the curl cmd
type CurlRequest struct {
	Host        string
	Path        string
	Port string
	QueryParams map[string]interface{}
	Headers     map[string]interface{}
}
