package models

type apiErrors []struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type apiResult struct {
	Ipv4Cidrs []string `json:"ipv4_cidrs"`
	Ipv6Cidrs []string `json:"ipv6_cidrs"`
	Etag      string   `json:"etag"`
}

type ApiResponse struct {
	Result   apiResult `json:"result"`
	Success  bool      `json:"success"`
	Errors   apiErrors `json:"errors"`
	Messages []string  `json:"messages"`
}
