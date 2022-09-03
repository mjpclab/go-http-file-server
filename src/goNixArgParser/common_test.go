package goNixArgParser

func expectStrings(actuals []string, expects ...string) bool {
	if len(actuals) != len(expects) {
		return false
	}

	for i := range actuals {
		if actuals[i] != expects[i] {
			return false
		}
	}

	return true
}
