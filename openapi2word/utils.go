package openapi2word

import "github.com/go-courier/oas"

func CheckMethod(method oas.HttpMethod) string {
	switch method {
	case "get":
		return "GET"
	case "post":
		return "POST"
	case "put":
		return "PUT"
	case "delete":
		return "DELETE"
	}
	return string(method)
}

func CheckType(t string) string {
	switch t {
	case "GitQuerycapComToolsDatatypesUUID":
		return "string"
	case "GithubComGoCourierSqlxV2DatatypesMySQLTimestamp":
		return "time"
	case "GithubComGoCourierSqlxV2DatatypesBool":
		return "Bool"
	}

	if t == "" {
		t = "string"
	}
	return t
}
