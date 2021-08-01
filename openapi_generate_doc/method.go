package openapi_generate_doc

import "github.com/go-courier/oas"

func Method(method oas.HttpMethod) string {
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
