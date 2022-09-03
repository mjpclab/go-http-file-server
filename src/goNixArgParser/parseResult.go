package goNixArgParser

// =============================
// set configOptions
// =============================

func (r *ParseResult) SetConfigOption(key, value string) {
	r.configOptions[key] = []string{value}
}

func (r *ParseResult) SetConfigOptions(key string, values []string) {
	var configValues []string

	if opt := r.keyOptionMap[key]; opt != nil {
		configValues = opt.filterValues(values)
	} else {
		configValues = copys(values)
	}

	r.configOptions[key] = configValues
}

//=============================
// has xxx
//=============================

func (r *ParseResult) HasFlagKey(key string) bool {
	_, found := r.specifiedOptions[key]
	return found
}

func (r *ParseResult) HasFlagValue(key string) bool {
	return len(r.specifiedOptions[key]) > 0
}

func (r *ParseResult) HasEnvKey(key string) bool {
	_, found := r.envs[key]
	return found
}

func (r *ParseResult) HasEnvValue(key string) bool {
	return len(r.envs[key]) > 0
}

func (r *ParseResult) HasConfigKey(key string) bool {
	_, found := r.configOptions[key]
	return found
}

func (r *ParseResult) HasConfigValue(key string) bool {
	return len(r.configOptions[key]) > 0
}

func (r *ParseResult) HasDefaultKey(key string) bool {
	_, found := r.defaults[key]
	return found
}

func (r *ParseResult) HasDefaultValue(key string) bool {
	return len(r.defaults[key]) > 0
}

func (r *ParseResult) HasKey(key string) bool {
	return r.HasFlagKey(key) || r.HasEnvKey(key) || r.HasConfigKey(key) || r.HasDefaultKey(key)
}

func (r *ParseResult) HasValue(key string) bool {
	return r.HasFlagValue(key) || r.HasEnvValue(key) || r.HasConfigValue(key) || r.HasDefaultValue(key)
}

//=============================
// get single value
//=============================

func (r *ParseResult) GetString(key string) (value string, found bool) {
	value, found = getValue(r.specifiedOptions, key)
	if found {
		return
	}

	value, found = getValue(r.envs, key)
	if found {
		return
	}

	value, found = getValue(r.configOptions, key)
	if found {
		return
	}

	value, found = getValue(r.defaults, key)
	if found {
		return
	}

	return
}

func (r *ParseResult) GetBool(key string) (value bool, found bool) {
	str, found := r.GetString(key)
	if !found {
		return
	}

	value, err := toBool(str)
	found = err == nil
	return
}

func (r *ParseResult) GetInt(key string) (value int, found bool) {
	str, found := r.GetString(key)
	if !found {
		return
	}

	value, err := toInt(str)
	found = err == nil
	return
}

func (r *ParseResult) GetInt64(key string) (value int64, found bool) {
	str, found := r.GetString(key)
	if !found {
		return
	}

	value, err := toInt64(str)
	found = err == nil
	return
}

func (r *ParseResult) GetUint64(key string) (value uint64, found bool) {
	str, found := r.GetString(key)
	if !found {
		return
	}

	value, err := toUint64(str)
	found = err == nil
	return
}

func (r *ParseResult) GetFloat64(key string) (value float64, found bool) {
	str, found := r.GetString(key)
	if !found {
		return
	}

	value, err := toFloat64(str)
	found = err == nil
	return
}

// =============================
// get multi values
// =============================

func (r *ParseResult) GetStrings(key string) (values []string, found bool) {
	values, found = getValues(r.specifiedOptions, key)
	if found {
		return
	}

	values, found = getValues(r.envs, key)
	if found {
		return
	}

	values, found = getValues(r.configOptions, key)
	if found {
		return
	}

	values, found = getValues(r.defaults, key)
	if found {
		return
	}

	return
}

func (r *ParseResult) GetBools(key string) (values []bool, found bool) {
	strs, found := r.GetStrings(key)
	if !found {
		return
	}

	values, err := toBools(strs)
	found = err == nil
	return
}

func (r *ParseResult) GetInts(key string) (values []int, found bool) {
	strs, found := r.GetStrings(key)
	if !found {
		return
	}

	values, err := toInts(strs)
	found = err == nil
	return
}

func (r *ParseResult) GetInt64s(key string) (values []int64, found bool) {
	strs, found := r.GetStrings(key)
	if !found {
		return
	}

	values, err := toInt64s(strs)
	found = err == nil
	return
}

func (r *ParseResult) GetUint64s(key string) (values []uint64, found bool) {
	strs, found := r.GetStrings(key)
	if !found {
		return
	}

	values, err := toUint64s(strs)
	found = err == nil
	return
}

func (r *ParseResult) GetFloat64s(key string) (values []float64, found bool) {
	strs, found := r.GetStrings(key)
	if !found {
		return
	}

	values, err := toFloat64s(strs)
	found = err == nil
	return
}

func (r *ParseResult) GetRests() (rests []string) {
	if len(r.specifiedRests) > 0 {
		return copys(r.specifiedRests)
	} else if len(r.configRests) > 0 {
		return copys(r.configRests)
	}

	return
}

// =============================
// commands
// =============================

func (r *ParseResult) GetCommands() []string {
	return copys(r.commands)
}

// =============================
// ambigus
// =============================

func (r *ParseResult) HasAmbigu() bool {
	return len(r.specifiedAmbigus) > 0 || len(r.configAmbigus) > 0
}

func (r *ParseResult) GetAmbigus() []string {
	flags := make([]string, 0, len(r.specifiedAmbigus)+len(r.configAmbigus))

	for _, flag := range r.specifiedAmbigus {
		if !contains(flags, flag) {
			flags = append(flags, flag)
		}
	}

	for _, flag := range r.configAmbigus {
		if !contains(flags, flag) {
			flags = append(flags, flag)
		}
	}

	return flags
}

// =============================
// undefs
// =============================

func (r *ParseResult) HasUndef() bool {
	return len(r.specifiedUndefs) > 0 || len(r.configUndefs) > 0
}

func (r *ParseResult) GetUndefs() []string {
	flags := make([]string, 0, len(r.specifiedUndefs)+len(r.configUndefs))

	for _, flag := range r.specifiedUndefs {
		if !contains(flags, flag) {
			flags = append(flags, flag)
		}
	}

	for _, flag := range r.configUndefs {
		if !contains(flags, flag) {
			flags = append(flags, flag)
		}
	}

	return flags
}
