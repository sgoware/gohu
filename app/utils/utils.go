package utils

import (
	"strings"
)

func GetServiceFullName(serviceName string) string {
	output := "gohu-"
	for _, v := range serviceName {
		if IsUpperCase(v) {
			output += "-" + string(ToLower(v))
		} else if v == '.' {
			output += "-"
		} else {
			output += string(v)
		}
	}
	return output
}

func GetNamespace(serviceName string) string {
	output := strings.Split(serviceName, ".")
	return output[0] + ".yaml"
}

func GetNamespaceType(namespace string) string {
	output := strings.Split(namespace, ".")
	if len(output) == 1 {
		return "properties"
	}
	return output[1]
}

func GetServiceType(serviceName string) string {
	output := strings.Split(serviceName, ".")
	return output[1]
}

func GetServiceSingleName(serviceName string) string {
	output := strings.Split(serviceName, ".")
	if len(output) == 2 {
		return output[0]
	}
	return output[2]
}

func GetServiceDetails(serviceName string) (string, string, string) {
	output := strings.Split(serviceName, ".")
	if len(output) == 2 {
		return output[0] + ".yaml", output[1], output[0]
	}
	return output[0] + ".yaml", output[1], output[2]
}

func IsUpperCase(char int32) bool {
	if char >= 65 && char <= 90 {
		return true
	}
	return false
}

func ToLower(char int32) int32 {
	if char >= 65 && char <= 90 {
		return char + 'a' - 'A'
	}
	return char
}
