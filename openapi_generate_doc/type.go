package openapi_generate_doc

func Type(t string) string {
	switch t {
	case "GitQuerycapComToolsDatatypesUUID",
		"GitQuerycapComToolsDatatypesUUIDs":
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
