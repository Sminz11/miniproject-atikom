package util

const MockAuthCode = "mock-auth-code"
const MockToken = "mock-jwt-token"

func ValidateCode(code string) bool {
	return code == MockAuthCode
}
