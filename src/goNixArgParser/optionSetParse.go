package goNixArgParser

import (
	"strings"
)

func (s *OptionSet) splitMergedArg(arg *Arg) (args []*Arg, success bool) {
	flagMap := s.flagMap
	optionMap := s.flagOptionMap
	argText := arg.Text

	if arg.Type != UndetermArg ||
		len(argText) <= len(s.mergeFlagPrefix) ||
		!strings.HasPrefix(argText, s.mergeFlagPrefix) {
		return
	}

	if flagMap[argText] != nil {
		return
	}

	var prevFlag *Flag
	mergedArgs := argText[len(s.mergeFlagPrefix):]
	splittedArgs := make([]*Arg, 0, len(mergedArgs))
	for i, mergedArg := range mergedArgs {
		splittedArg := s.mergeFlagPrefix + string(mergedArg)
		flag := flagMap[splittedArg]

		if flag != nil {
			if !flag.canMerge {
				return
			}
			splittedArgs = append(splittedArgs, NewArg(splittedArg, FlagArg))
			prevFlag = flag
			continue
		}

		if len(splittedArg) <= 1 {
			return
		}

		if prevFlag == nil {
			return
		}

		option := optionMap[prevFlag.Name]
		if option == nil || !option.AcceptValue {
			return
		}

		// re-generate standalone flag with values
		splittedArgs[len(splittedArgs)-1] = NewArg(prevFlag.Name+mergedArgs[i:], UndetermArg)
		break
	}

	return splittedArgs, true
}

func (s *OptionSet) splitMergedArgs(initArgs []*Arg) []*Arg {
	args := make([]*Arg, 0, len(initArgs))
	for _, arg := range initArgs {
		splittedArgs, splitted := s.splitMergedArg(arg)
		if splitted {
			args = append(args, splittedArgs...)
		} else {
			args = append(args, arg)
		}
	}
	return args
}

func (s *OptionSet) splitAssignSignArg(arg *Arg) (args []*Arg) {
	args = make([]*Arg, 0, 2)

	if arg.Type != UndetermArg {
		args = append(args, arg)
		return
	}

	argText := arg.Text
	for _, flag := range s.flagMap {
		flagName := flag.Name
		if !s.flagOptionMap[flagName].AcceptValue {
			continue
		}
		for _, assignSign := range flag.assignSigns {
			if len(assignSign) == 0 {
				continue
			}

			prefix := flagName + assignSign
			if strings.HasPrefix(argText, prefix) {
				args = append(args, NewArg(flagName, FlagArg))
				args = append(args, NewArg(argText[len(prefix):], ValueArg))
				return
			}

			assignIndex := strings.Index(argText, assignSign)
			if assignIndex <= 0 {
				continue
			}
			prefix = argText[0:assignIndex]
			if foundFlag, _ := s.findFlagByPrefix(prefix); foundFlag == flag {
				args = append(args, NewArg(flagName, FlagArg))
				args = append(args, NewArg(argText[assignIndex+len(assignSign):], ValueArg))
				return
			}
		}
	}

	args = append(args, arg)
	return
}

func (s *OptionSet) splitAssignSignArgs(initArgs []*Arg) []*Arg {
	args := make([]*Arg, 0, len(initArgs))

	for _, initArg := range initArgs {
		args = append(args, s.splitAssignSignArg(initArg)...)
	}

	return args
}

func (s *OptionSet) splitConcatAssignArg(arg *Arg) (args []*Arg) {
	args = make([]*Arg, 0, 2)

	if arg.Type != UndetermArg {
		args = append(args, arg)
		return
	}

	argText := arg.Text
	for _, flag := range s.flagMap {
		if !flag.canConcatAssign ||
			!s.flagOptionMap[flag.Name].AcceptValue ||
			len(argText) <= len(flag.Name) ||
			!strings.HasPrefix(argText, flag.Name) {
			continue
		}
		flagName := flag.Name
		flagValue := argText[len(flagName):]
		args = append(args, NewArg(flagName, FlagArg))
		args = append(args, NewArg(flagValue, ValueArg))
		return
	}

	args = append(args, arg)
	return
}

func (s *OptionSet) splitConcatAssignArgs(initArgs []*Arg) []*Arg {
	args := make([]*Arg, 0, len(initArgs))

	for _, initArg := range initArgs {
		args = append(args, s.splitConcatAssignArg(initArg)...)
	}

	return args
}

func (s *OptionSet) markAmbiguPrefixArgsValues(args []*Arg) {
	foundAmbiguFlag := false
	for _, arg := range args {
		if arg.Type != UndetermArg {
			foundAmbiguFlag = false
			continue
		}
		actualFlag, ambiguous := s.findFlagByPrefix(arg.Text)
		if ambiguous {
			arg.Type = AmbiguousFlagArg
			foundAmbiguFlag = true
		} else if actualFlag != nil {
			arg.Type = FlagArg
			arg.Text = actualFlag.Name
			foundAmbiguFlag = false
		} else if foundAmbiguFlag {
			arg.Type = AmbiguousFlagValueArg
		}
	}
}

func (s *OptionSet) markUndefArgsValues(args []*Arg) {
	foundUndefFlag := false
	for _, arg := range args {
		if arg.Type != UndetermArg {
			foundUndefFlag = false
			continue
		}
		if s.isUdefFlag(arg.Text) {
			arg.Type = UndefFlagArg
			foundUndefFlag = true
		} else if foundUndefFlag {
			arg.Type = UndefFlagValueArg
		}
	}
}

func isValueArg(flag *Flag, arg *Arg) bool {
	switch arg.Type {
	case ValueArg:
		return true
	case UndetermArg:
		return flag.canFollowAssign
	default:
		return false
	}
}

func (s *OptionSet) parseArgsInGroup(argObjs []*Arg) (args map[string][]string, rests, ambigus, undefs []string) {
	args = map[string][]string{}
	rests = []string{}
	ambigus = []string{}
	undefs = []string{}

	flagOptionMap := s.flagOptionMap
	flagMap := s.flagMap

	if s.hasCanMerge {
		argObjs = s.splitMergedArgs(argObjs)
	}
	if s.hasAssignSigns {
		argObjs = s.splitAssignSignArgs(argObjs)
	}
	if s.hasCanConcatAssign {
		argObjs = s.splitConcatAssignArgs(argObjs)
	}

	s.markAmbiguPrefixArgsValues(argObjs)
	s.markUndefArgsValues(argObjs)

	// walk
	for i, argCount, peeked := 0, len(argObjs), 0; i < argCount; i, peeked = i+1+peeked, 0 {
		arg := argObjs[i]

		// rests
		if arg.Type == RestSignArg {
			continue
		}

		if arg.Type == UndetermArg {
			arg.Type = RestArg
		}
		if arg.Type == RestArg {
			rests = append(rests, arg.Text)
			continue
		}

		// ambigus
		if arg.Type == AmbiguousFlagValueArg {
			continue
		}

		if arg.Type == AmbiguousFlagArg {
			ambigus = append(ambigus, arg.Text)
			continue
		}

		// undefs
		if arg.Type == UndefFlagValueArg {
			continue
		}

		if arg.Type == UndefFlagArg {
			undefs = append(undefs, arg.Text)
			continue
		}

		// normal
		opt := flagOptionMap[arg.Text]
		flag := flagMap[arg.Text]

		if !opt.AcceptValue { // option has no value
			args[opt.Key] = []string{}
			continue
		}

		if !opt.MultiValues { // option has 1 value
			if i == argCount-1 || !isValueArg(flag, argObjs[i+1]) { // no more value
				if opt.OverridePrev || args[opt.Key] == nil {
					args[opt.Key] = []string{}
				}
			} else {
				if opt.OverridePrev || args[opt.Key] == nil {
					nextArg := argObjs[i+1]
					nextArg.Type = ValueArg
					args[opt.Key] = []string{nextArg.Text}
				}
				peeked++
			}
			continue
		}

		//option have multi values
		values := []string{}
		for {
			if i+peeked == argCount-1 { // last arg reached
				break
			}

			if !isValueArg(flag, argObjs[i+peeked+1]) { // no more value
				break
			}

			peeked++
			peekedArg := argObjs[i+peeked]
			peekedArg.Type = ValueArg
			value := peekedArg.Text
			var appending []string
			if len(opt.Delimiters) == 0 {
				appending = []string{value}
			} else {
				appending = strings.FieldsFunc(value, opt.isDelimiter)
			}

			if opt.UniqueValues {
				values = appendUnique(values, appending...)
			} else {
				values = append(values, appending...)
			}
		}

		if opt.OverridePrev || args[opt.Key] == nil {
			args[opt.Key] = values
		} else {
			args[opt.Key] = append(args[opt.Key], values...)
		}
	}

	return args, rests, ambigus, undefs
}

func (s *OptionSet) parseInGroup(argObjs, configObjs []*Arg) *ParseResult {
	keyOptionMap := s.keyOptionMap

	args, argRests, argAmbigus, argUndefs := s.parseArgsInGroup(argObjs)
	envs := s.keyEnvMap
	configs, configRests, configAmbigus, configUndefs := s.parseArgsInGroup(configObjs)
	defaults := s.keyDefaultMap

	return &ParseResult{
		keyOptionMap: keyOptionMap,

		args:     args,
		envs:     envs,
		configs:  configs,
		defaults: defaults,

		argRests:    argRests,
		configRests: configRests,

		argAmbigus:    argAmbigus,
		configAmbigus: configAmbigus,

		argUndefs:    argUndefs,
		configUndefs: configUndefs,
	}
}

func (s *OptionSet) getNormalizedArgs(initArgs []string) []*Arg {
	args := make([]*Arg, 0, len(initArgs)+1)

	foundRestSign := false
	for _, arg := range initArgs {
		switch {
		case s.isGroupSep(arg):
			foundRestSign = false
			args = append(args, NewArg(arg, GroupSepArg))
		case foundRestSign:
			args = append(args, NewArg(arg, RestArg))
		case s.isRestSign(arg):
			foundRestSign = true
			args = append(args, NewArg(arg, RestSignArg))
		case s.flagMap[arg] != nil:
			args = append(args, NewArg(arg, FlagArg))
		default:
			args = append(args, NewArg(arg, UndetermArg))
		}
	}

	return args
}

func splitArgsIntoGroups(argObjs []*Arg) [][]*Arg {
	argObjs = append(argObjs, NewArg("", GroupSepArg))

	groups := [][]*Arg{}
	items := []*Arg{}
	for _, argObj := range argObjs {
		if argObj.Type != GroupSepArg {
			items = append(items, argObj)
			continue
		}

		groups = append(groups, items)
		items = []*Arg{}
	}

	return groups
}

func (s *OptionSet) getArgsConfigsGroups(initArgs, initConfigs []string) ([][]*Arg, [][]*Arg) {
	args := s.getNormalizedArgs(initArgs)
	argsGroups := splitArgsIntoGroups(args)
	argsGroupsCount := len(argsGroups)

	configs := s.getNormalizedArgs(initConfigs)
	configsGroups := splitArgsIntoGroups(configs)
	configsGroupsCount := len(configsGroups)

	length := argsGroupsCount
	if configsGroupsCount > length {
		length = configsGroupsCount
	}

	for i := 0; i < length-argsGroupsCount; i++ {
		argsGroups = append(argsGroups, []*Arg{})
	}

	for i := 0; i < length-configsGroupsCount; i++ {
		configsGroups = append(configsGroups, []*Arg{})
	}

	return argsGroups, configsGroups
}

func (s *OptionSet) ParseGroups(initArgs, initConfigs []string) []*ParseResult {
	argsGroups, configsGroups := s.getArgsConfigsGroups(initArgs, initConfigs)

	results := []*ParseResult{}
	for i, length := 0, len(argsGroups); i < length; i++ {
		result := s.parseInGroup(argsGroups[i], configsGroups[i])
		results = append(results, result)
	}

	return results
}

func (s *OptionSet) Parse(initArgs, initConfigs []string) *ParseResult {
	argsGroups, configsGroups := s.getArgsConfigsGroups(initArgs, initConfigs)

	var args []*Arg
	if len(argsGroups) > 0 {
		args = argsGroups[0]
	} else {
		args = []*Arg{}
	}

	var configs []*Arg
	if len(configsGroups) > 0 {
		configs = configsGroups[0]
	} else {
		configs = []*Arg{}
	}

	result := s.parseInGroup(args, configs)

	return result
}
