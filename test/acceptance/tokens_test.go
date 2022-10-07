package acceptance

func tstNoToken() string {
	return ""
}

func tstValidApiToken() string {
	return "put_secure_random_string_here_for_api_token_test_token"
}

func tstInvalidApiToken() string {
	return "invalid_put_secure_random_string_here_for_api_token_test_token"
}
