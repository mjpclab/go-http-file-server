package goNixArgParser

func (r *ParseResult) HasFlagKey(key string) bool {
	_, found := r.params[key]
	return found
}

func (r *ParseResult) HasFlagValue(key string) bool {
	return len(r.params[key]) > 0
}

func (r *ParseResult) HasEnvKey(key string) bool {
	_, found := r.envs[key]
	return found
}

func (r *ParseResult) HasEnvValue(key string) bool {
	return len(r.envs[key]) > 0
}

func (r *ParseResult) HasDefaultKey(key string) bool {
	_, found := r.defaults[key]
	return found
}

func (r *ParseResult) HasDefaultValue(key string) bool {
	return len(r.defaults[key]) > 0
}

func _getValue(source map[string][]string, key string) (value string, found bool) {
	var values []string
	values, found = source[key]

	if found && len(values) > 0 {
		value = values[0]
	}

	return
}

func (r *ParseResult) GetValue(key string) (value string, found bool) {
	value, found = _getValue(r.params, key)

	if !found {
		value, found = _getValue(r.envs, key)
	}

	if !found {
		value, found = _getValue(r.defaults, key)
	}

	return
}

func _getValues(source map[string][]string, key string) (values []string, found bool) {
	sourceValues, found := source[key]
	if found {
		values = make([]string, len(sourceValues))
		copy(values, sourceValues)
		return values, true
	}
	return
}

func (r *ParseResult) GetValues(key string) (values []string, found bool) {
	values, found = _getValues(r.params, key)

	if !found {
		values, found = _getValues(r.envs, key)
	}

	if !found {
		values, found = _getValues(r.defaults, key)
	}

	return
}

func (r *ParseResult) GetRests() []string {
	rests := make([]string, len(r.rests))
	copy(rests, r.rests)
	return rests
}
