package utils

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"strings"
)

func JoinStrings(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

func GenerateID(prefix string) string {
	nanoid, err := gonanoid.Generate(defaultAlphabet, defaultSize)
	if err != nil {
		return ""
	}

	return prefix + "_" + nanoid
}

func GetErrorMetaData(err error) map[string]string {
	md := make(map[string]string)
	md["err"] = err.Error()

	return md
}
