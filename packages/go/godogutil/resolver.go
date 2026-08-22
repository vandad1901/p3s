//go:build test

package godogutil

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	randomStringRegex = regexp.MustCompile(`\{([0-9]+)(c|d)?\}`)
	referenceRegex    = regexp.MustCompile(`\{([\$a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)\}`)
)

const emptyPattern = "{$empty}"

func ResolveReference(s *SharedData, row map[string]string, key string) string {
	return referenceRegex.ReplaceAllStringFunc(row[key], func(match string) string {
		subMatches := referenceRegex.FindStringSubmatch(match)
		rowKey := subMatches[1]
		fieldKey := subMatches[2]

		row, exists := s.DataMap[rowKey]
		s.Require.True(exists, "row not found for reference: "+rowKey)

		value, exists := row[fieldKey]
		s.Require.True(exists, "field not found in row for reference: "+fieldKey)

		return value
	})
}

func ResolveString(s *SharedData, row map[string]string, key string) string {
	valueStr := ResolveReference(s, row, key)
	if valueStr == "" || valueStr == emptyPattern {
		valueStr = ""
	}

	value := randomStringRegex.ReplaceAllStringFunc(valueStr, func(match string) string {
		subMatches := randomStringRegex.FindStringSubmatch(match)
		length, err := strconv.Atoi(subMatches[1])
		s.Require.NoError(err, "failed to convert length to int for key: "+key)

		switch subMatches[2] {
		case "c":
			return generateRandomChars(length)
		case "d":
			return generateRandomNumbers(length)
		case "":
			return generateRandomCharNumbers(length)
		default:
			s.Require.Fail("invalid random type for key: " + key)
		}

		return ""
	})

	s.DataMap[row["Key"]][key] = value

	return value
}

func ResolveInt64(s *SharedData, row map[string]string, key string) int64 {
	valueStr := ResolveReference(s, row, key)
	if valueStr == "" || valueStr == emptyPattern {
		valueStr = "0"
	}

	s.DataMap[row["Key"]][key] = valueStr
	value, err := strconv.ParseInt(valueStr, 10, 64)
	s.Require.NoError(err, "failed to convert value to int64 for key: "+key)

	return value
}

func ResolveInt32(s *SharedData, row map[string]string, key string) int32 {
	valueStr := ResolveReference(s, row, key)
	if valueStr == "" || valueStr == emptyPattern {
		valueStr = "0"
	}

	value, err := strconv.ParseInt(valueStr, 10, 32)
	s.Require.NoError(err, "failed to convert value to int32 for key: "+key)

	s.DataMap[row["Key"]][key] = valueStr

	return int32(value)
}

func ResolveBoolean(s *SharedData, row map[string]string, key string) bool {
	valueStr := ResolveReference(s, row, key)
	if valueStr == "" || valueStr == emptyPattern {
		valueStr = "false"
	}

	value, err := strconv.ParseBool(valueStr)
	s.Require.NoError(err, "failed to convert value to bool for key: "+key)

	s.DataMap[row["Key"]][key] = valueStr

	return value
}

func ResolveEnumFactory(s *SharedData, row map[string]string, key string, enumValues []string) int32 {
	valueStr := ResolveReference(s, row, key)
	if valueStr == "" || valueStr == emptyPattern {
		valueStr = "0"
	}

	valueStr = strings.ToLower(valueStr)

	for i, enumValue := range enumValues {
		indexStr := strconv.Itoa(i)
		if "$"+enumValue == valueStr || indexStr == valueStr {
			s.DataMap[row["Key"]][key] = indexStr

			return int32(i)
		}
	}

	s.Require.Fail("invalid enum value for key: " + key + ", value: " + valueStr)

	return 0
}

func commaSeparatedValueToSlice(value string) []string {
	if value == "" {
		return []string{}
	}

	return strings.Split(value, ",")
}

func ResolveCommaSeparatedValues(s *SharedData, row map[string]string, key string) []string {
	valueStr := ResolveReference(s, row, key)
	if valueStr == "" || valueStr == emptyPattern {
		valueStr = ""
	}

	s.DataMap[row["Key"]][key] = valueStr

	return commaSeparatedValueToSlice(valueStr)
}

const TimeFMT = time.RFC3339Nano

func ResolveTime(s *SharedData, row map[string]string, key string) time.Time {
	valueStr := ResolveReference(s, row, key)
	if valueStr == "" || valueStr == emptyPattern {
		valueStr = "0001-01-01T00:00:00Z"
	}

	value, err := time.Parse(TimeFMT, valueStr)
	s.Require.NoError(err)

	s.DataMap[row["Key"]][key] = valueStr

	return value
}
