


```go
func AzureADEndpoint(tenant string) oauth2.Endpoint {
	if tenant == "" {
		tenant = "common"
	}
	return oauth2.Endpoint{
		AuthURL:       "https://login.microsoftonline.com/" + tenant + "/oauth2/v2.0/authorize",
		TokenURL:      "https://login.microsoftonline.com/" + tenant + "/oauth2/v2.0/token",
		DeviceAuthURL: "https://login.microsoftonline.com/" + tenant + "/oauth2/v2.0/devicecode",
	}
}
```