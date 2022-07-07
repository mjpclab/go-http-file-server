package goNixArgParser

import (
	"strings"
)

func (s *OptionSet) splitMergedToken(token *argToken) (results []*argToken, success bool) {
	flagMap := s.nameFlagMap
	optionMap := s.flagOptionMap
	originalArg := token.text

	if token.kind != undetermArg ||
		len(originalArg) <= len(s.mergeFlagPrefix) ||
		!strings.HasPrefix(originalArg, s.mergeFlagPrefix) {
		return
	}

	if flagMap[originalArg] != nil {
		return
	}

	var prevFlag *Flag
	mergedArgs := originalArg[len(s.mergeFlagPrefix):]
	splittedTokens := make([]*argToken, 0, len(mergedArgs))
	for i, mergedArg := range mergedArgs {
		splittedArg := s.mergeFlagPrefix + string(mergedArg)
		flag := flagMap[splittedArg]

		if flag != nil {
			if !flag.canMerge {
				return
			}
			splittedTokens = append(splittedTokens, newToken(splittedArg, flagArg))
			prevFlag = flag
			continue
		}

		if prevFlag == nil {
			return
		}

		option := optionMap[prevFlag.Name]
		if option == nil || !option.AcceptValue {
			return
		}

		// re-generate standalone flag with values
		splittedTokens[len(splittedTokens)-1] = newToken(prevFlag.Name+mergedArgs[i:], undetermArg)
		break
	}

	return splittedTokens, true
}

func (s *OptionSet) splitMergedTokens(tokens []*argToken) []*argToken {
	results := make([]*argToken, 0, len(tokens))
	for _, originalToken := range tokens {
		splittedTokens, splitted := s.splitMergedToken(originalToken)
		if splitted {
			results = append(results, splittedTokens...)
		} else {
			results = append(results, originalToken)
		}
	}
	return results
}

func (s *OptionSet) splitAssignSignToken(token *argToken) (results []*argToken) {
	results = make([]*argToken, 0, 2)

	text := token.text
	for _, flag := range s.nameFlagMap {
		flagName := flag.Name
		if !s.flagOptionMap[flagName].AcceptValue {
			continue
		}
		for _, assignSign := range s.assignSigns {
			if len(assignSign) == 0 {
				continue
			}

			prefix := flagName + assignSign
			if strings.HasPrefix(text, prefix) {
				results = append(results,
					newToken(flagName, flagArg),
					newToken(text[len(prefix):], valueArg),
				)
				return
			}

			assignIndex := strings.Index(text, assignSign)
			if assignIndex <= 0 {
				continue
			}
			prefix = text[0:assignIndex]
			if foundFlag, _ := s.findFlagByPrefix(prefix); foundFlag == flag {
				results = append(results,
					newToken(flagName, flagArg),
					newToken(text[assignIndex+len(assignSign):], valueArg),
				)
				return
			}
		}
	}

	results = append(results, token)
	return
}

func (s *OptionSet) splitAssignSignTokens(tokens []*argToken) []*argToken {
	results := make([]*argToken, 0, len(tokens))

	for _, token := range tokens {
		if token.kind == undetermArg {
			results = append(results, s.splitAssignSignToken(token)...)
		} else {
			results = append(results, token)
		}
	}

	return results
}

func (s *OptionSet) splitConcatAssignToken(token *argToken) (results []*argToken) {
	results = make([]*argToken, 0, 2)

	text := token.text
	for _, flag := range s.nameFlagMap {
		if !flag.canConcatAssign ||
			!s.flagOptionMap[flag.Name].AcceptValue ||
			len(text) <= len(flag.Name) ||
			!strings.HasPrefix(text, flag.Name) {
			continue
		}
		flagName := flag.Name
		flagValue := text[len(flagName):]
		results = append(results,
			newToken(flagName, flagArg),
			newToken(flagValue, valueArg),
		)
		return
	}

	results = append(results, token)
	return
}

func (s *OptionSet) splitConcatAssignTokens(tokens []*argToken) []*argToken {
	results := make([]*argToken, 0, len(tokens))

	for _, token := range tokens {
		if token.kind == undetermArg {
			results = append(results, s.splitConcatAssignToken(token)...)
		} else {
			results = append(results, token)
		}
	}

	return results
}

func (s *OptionSet) markAmbiguPrefixTokens(tokens []*argToken) {
	foundAmbiguFlag := false
	for _, token := range tokens {
		if token.kind != undetermArg {
			foundAmbiguFlag = false
			continue
		}
		actualFlag, ambiguous := s.findFlagByPrefix(token.text)
		if ambiguous {
			token.kind = ambiguousFlagArg
			foundAmbiguFlag = true
		} else if actualFlag != nil {
			token.kind = flagArg
			token.text = actualFlag.Name
			foundAmbiguFlag = false
		} else if foundAmbiguFlag {
			token.kind = ambiguousFlagValueArg
		}
	}
}

func (s *OptionSet) markUndefTokens(tokens []*argToken) {
	for _, token := range tokens {
		if token.kind != undetermArg {
			continue
		}
		if s.isUdefFlag(token.text) {
			token.kind = undefFlagArg
			// remove assign value
			for _, assignSign := range s.assignSigns {
				assignIndex := strings.Index(token.text, assignSign)
				if assignIndex > 0 {
					token.text = token.text[:assignIndex]
					break
				}
			}
		}
	}
}

func isValueToken(flag *Flag, token *argToken) bool {
	switch token.kind {
	case valueArg:
		return true
	case undetermArg:
		return flag.canFollowAssign
	default:
		return false
	}
}

func (s *OptionSet) parseTokensInGroup(tokens []*argToken) (options map[string][]string, rests, ambigus, undefs []string) {
	options = map[string][]string{}
	rests = []string{}
	ambigus = []string{}
	undefs = []string{}

	flagOptionMap := s.flagOptionMap
	flagMap := s.nameFlagMap

	if s.hasCanMerge {
		tokens = s.splitMergedTokens(tokens)
	}
	if len(s.assignSigns) > 0 {
		tokens = s.splitAssignSignTokens(tokens)
	}
	if s.hasCanConcatAssign {
		tokens = s.splitConcatAssignTokens(tokens)
	}

	s.markAmbiguPrefixTokens(tokens)
	s.markUndefTokens(tokens)

	// walk
	for i, tokenCount, peeked := 0, len(tokens), 0; i < tokenCount; i, peeked = i+1+peeked, 0 {
		token := tokens[i]

		// rests
		if token.kind == restSignArg {
			continue
		}

		if token.kind == undetermArg {
			token.kind = restArg
		}
		if token.kind == restArg {
			rests = append(rests, token.text)
			continue
		}

		// ambigus
		if token.kind == ambiguousFlagValueArg {
			continue
		}

		if token.kind == ambiguousFlagArg {
			ambigus = append(ambigus, token.text)
			continue
		}

		// undefs
		if token.kind == undefFlagValueArg {
			continue
		}

		if token.kind == undefFlagArg {
			undefs = append(undefs, token.text)
			continue
		}

		// normal
		opt := flagOptionMap[token.text]
		flag := flagMap[token.text]

		if !opt.AcceptValue { // option has no value
			options[opt.Key] = []string{}
			continue
		}

		if !opt.MultiValues { // option has 1 value
			if i == tokenCount-1 || !isValueToken(flag, tokens[i+1]) { // no more value
				if opt.OverridePrev || options[opt.Key] == nil {
					options[opt.Key] = []string{}
				}
			} else {
				if opt.OverridePrev || options[opt.Key] == nil {
					nextArg := tokens[i+1]
					nextArg.kind = valueArg
					options[opt.Key] = []string{nextArg.text}
				}
				peeked++
			}
			continue
		}

		//option have multi values
		values := []string{}
		for {
			if i+peeked == tokenCount-1 { // last token reached
				break
			}

			if !isValueToken(flag, tokens[i+peeked+1]) { // no more value
				break
			}

			peeked++
			peekedToken := tokens[i+peeked]
			peekedToken.kind = valueArg
			text := peekedToken.text
			var appending []string
			if len(opt.Delimiters) == 0 {
				appending = []string{text}
			} else {
				appending = strings.FieldsFunc(text, opt.isDelimiter)
			}

			if opt.UniqueValues {
				values = appendUnique(values, appending...)
			} else {
				values = append(values, appending...)
			}
		}

		if opt.OverridePrev || options[opt.Key] == nil {
			options[opt.Key] = values
		} else {
			options[opt.Key] = append(options[opt.Key], values...)
		}
	}

	return options, rests, ambigus, undefs
}

func (s *OptionSet) parseInGroup(specifiedTokens, configTokens []*argToken) *ParseResult {
	keyOptionMap := s.keyOptionMap

	specifiedOptions, specifiedRests, specifiedAmbigus, specifiedUndefs := s.parseTokensInGroup(specifiedTokens)
	envs := s.keyEnvMap
	configOptions, configRests, configAmbigus, configUndefs := s.parseTokensInGroup(configTokens)
	defaults := s.keyDefaultMap

	return &ParseResult{
		keyOptionMap: keyOptionMap,

		specifiedOptions: specifiedOptions,
		envs:             envs,
		configOptions:    configOptions,
		defaults:         defaults,

		specifiedRests: specifiedRests,
		configRests:    configRests,

		specifiedAmbigus: specifiedAmbigus,
		configAmbigus:    configAmbigus,

		specifiedUndefs: specifiedUndefs,
		configUndefs:    configUndefs,
	}
}

func (s *OptionSet) argsToTokensGroups(args []string) (tokensGroups [][]*argToken) {
	tokensGroups = make([][]*argToken, 1)
	groupIndex := 0

	foundRestSign := false
	for _, arg := range args {
		switch {
		case s.isGroupSep(arg):
			tokensGroups = append(tokensGroups, make([]*argToken, 0, 4))
			groupIndex++
			foundRestSign = false
		case foundRestSign:
			tokensGroups[groupIndex] = append(tokensGroups[groupIndex], newToken(arg, restArg))
		case s.isRestSign(arg):
			tokensGroups[groupIndex] = append(tokensGroups[groupIndex], newToken(arg, restSignArg))
			foundRestSign = true
		case s.nameFlagMap[arg] != nil:
			tokensGroups[groupIndex] = append(tokensGroups[groupIndex], newToken(arg, flagArg))
		default:
			tokensGroups[groupIndex] = append(tokensGroups[groupIndex], newToken(arg, undetermArg))
		}
	}

	return
}

func (s *OptionSet) getAlignedTokensGroups(specifiedArgs, configArgs []string) ([][]*argToken, [][]*argToken) {
	specifiedTokensGroups := s.argsToTokensGroups(specifiedArgs)
	specifiedTokensGroupsCount := len(specifiedTokensGroups)

	configTokensGroups := s.argsToTokensGroups(configArgs)
	configTokensGroupsCount := len(configTokensGroups)

	maxCount := specifiedTokensGroupsCount
	if configTokensGroupsCount > maxCount {
		maxCount = configTokensGroupsCount
	}

	for i := specifiedTokensGroupsCount; i < maxCount; i++ {
		specifiedTokensGroups = append(specifiedTokensGroups, []*argToken{})
	}

	for i := configTokensGroupsCount; i < maxCount; i++ {
		configTokensGroups = append(configTokensGroups, []*argToken{})
	}

	return specifiedTokensGroups, configTokensGroups
}

func (s *OptionSet) ParseGroups(specifiedArgs, configArgs []string) []*ParseResult {
	specifiedTokensGroups, configTokensGroups := s.getAlignedTokensGroups(specifiedArgs, configArgs)

	length := len(specifiedTokensGroups)
	results := make([]*ParseResult, length)
	for i := 0; i < length; i++ {
		results[i] = s.parseInGroup(specifiedTokensGroups[i], configTokensGroups[i])
	}

	return results
}

func (s *OptionSet) Parse(specifiedArgs, configArgs []string) *ParseResult {
	specifiedTokensGroups, configTokensGroups := s.getAlignedTokensGroups(specifiedArgs, configArgs)

	var specifiedTokens []*argToken
	if len(specifiedTokensGroups) > 0 {
		specifiedTokens = specifiedTokensGroups[0]
	} else {
		specifiedTokens = []*argToken{}
	}

	var configTokens []*argToken
	if len(configTokensGroups) > 0 {
		configTokens = configTokensGroups[0]
	} else {
		configTokens = []*argToken{}
	}

	result := s.parseInGroup(specifiedTokens, configTokens)

	return result
}
