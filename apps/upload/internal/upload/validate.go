package upload

func validateMediaKey(key string) error {
	if key == "" {
		return errEmptyKey
	}

	for _, r := range key {
		if isValidCharacter(r) {
			continue
		}

		return errInvalidKey
	}

	return nil
}

func isValidCharacter(r rune) bool {
	if r >= 'a' && r <= 'z' ||
		r >= 'A' && r <= 'Z' ||
		r >= '0' && r <= '9' ||
		r == '-' || r == '_' {
		return true
	}

	return false
}
