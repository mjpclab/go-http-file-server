package goNixArgParser

func (r *ParseResult) HasKey(key string) bool {
	return r.params[key] != nil
}

func (r *ParseResult) HasValue(key string) bool {
	return r.params[key] != nil && len(r.params[key]) > 0
}

func _getValue(source map[string][]string, key string) (value string, found bool) {
	if source[key] != nil && len(source[key]) > 0 {
		return source[key][0], true
	}
	return
}

func (r *ParseResult) GetValue(key string) string {
	value, found := _getValue(r.params, key)
	if !found {
		value, _ = _getValue(r.defaults, key)
	}

	return value
}

func _getValues(source map[string][]string, key string) (values []string, found bool) {
	sourceValues := source[key]
	sourceValuesCount := len(sourceValues)
	if sourceValuesCount > 0 {
		values = make([]string, sourceValuesCount)
		copy(values, sourceValues)
		return values, true
	}
	return
}

func (r *ParseResult) GetValues(key string) (values []string) {
	values, found := _getValues(r.params, key)
	if !found {
		values, _ = _getValues(r.defaults, key)
	}

	return
}

func (r *ParseResult) GetDefaults(key string) []string {
	defaults := make([]string, len(r.defaults[key]))
	copy(defaults, r.defaults[key])
	return defaults
}

func (r *ParseResult) GetRests() []string {
	rests := make([]string, len(r.rests))
	copy(rests, r.rests)
	return rests
}
