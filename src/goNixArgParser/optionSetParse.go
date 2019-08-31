package goNixArgParser

import (
	"strings"
)

func (s *OptionSet) getNormalizedArgs(initArgs []string) []*Arg {
	args := make([]*Arg, 0, len(initArgs))

	for _, arg := range initArgs {
		if s.flagMap[arg] != nil {
			args = append(args, NewArg(arg, FlagArg))
		} else {
			args = append(args, NewArg(arg, UnknownArg))
		}
	}

	return args
}

func (s *OptionSet) splitMergedArg(arg *Arg) (args []*Arg, success bool) {
	flagMap := s.flagMap
	argText := arg.Text

	if arg.Type != UnknownArg ||
		len(argText) <= len(s.mergeOptionPrefix) ||
		!strings.HasPrefix(argText, s.mergeOptionPrefix) {
		return
	}

	mergedArgs := argText[len(s.mergeOptionPrefix):]
	splittedArgs := make([]*Arg, 0, len(mergedArgs))
	for _, mergedArg := range mergedArgs {
		splittedArg := s.mergeOptionPrefix + string(mergedArg)
		flag := flagMap[splittedArg]
		if flag == nil || !flag.canMerge {
			return
		}
		splittedArgs = append(splittedArgs, NewArg(splittedArg, FlagArg))
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

func (s *OptionSet) splitEqualAssignArg(arg *Arg) (args []*Arg) {
	args = make([]*Arg, 0, 2)

	if arg.Type != UnknownArg {
		args = append(args, arg)
		return
	}

	argText := arg.Text
	equalIndex := strings.IndexByte(argText, '=')
	if equalIndex == -1 {
		args = append(args, arg)
		return
	}

	flagName := argText[:equalIndex]
	flagValue := argText[equalIndex+1:]
	flag := s.flagMap[flagName]
	if flag == nil ||
		!flag.canEqualAssign ||
		!s.flagOptionMap[flagName].AcceptValue {
		args = append(args, arg)
		return
	}

	args = append(args, NewArg(flagName, FlagArg))
	args = append(args, NewArg(flagValue, ValueArg))
	return
}

func (s *OptionSet) splitEqualAssignArgs(initArgs []*Arg) []*Arg {
	args := make([]*Arg, 0, len(initArgs))

	for _, initArg := range initArgs {
		args = append(args, s.splitEqualAssignArg(initArg)...)
	}

	return args
}

func (s *OptionSet) splitConcatAssignArg(arg *Arg) (args []*Arg) {
	args = make([]*Arg, 0, 2)

	if arg.Type != UnknownArg {
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

func (s *OptionSet) Parse(initArgs []string) *ParseResult {
	params := map[string][]string{}
	defaults := s.keyDefaultMap
	rests := []string{}

	flagOptionMap := s.flagOptionMap

	args := s.getNormalizedArgs(initArgs)
	args = s.splitMergedArgs(args)
	args = s.splitEqualAssignArgs(args)
	args = s.splitConcatAssignArgs(args)

	// walk
	for i, argCount, peeked := 0, len(args), 0; i < argCount; i, peeked = i+1+peeked, 0 {
		arg := args[i]

		if arg.Type == UnknownArg {
			arg.Type = RestArg
			rests = append(rests, arg.Text)
			continue
		}

		opt := flagOptionMap[arg.Text]

		if !opt.AcceptValue { // option has no value
			params[opt.Key] = []string{}
			continue
		}

		if !opt.MultiValues { // option has 1 value
			if i == argCount-1 || args[i+1].Type == FlagArg { // no more value or next flag found
				if opt.OverridePrev || params[opt.Key] == nil {
					params[opt.Key] = []string{}
				}
			} else {
				if opt.OverridePrev || params[opt.Key] == nil {
					nextArg := args[i+1]
					nextArg.Type = ValueArg
					params[opt.Key] = []string{nextArg.Text}
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

			if args[i+peeked+1].Type == FlagArg { // next flag found
				break
			}

			peeked++
			peekedArg := args[i+peeked]
			peekedArg.Type = ValueArg
			value := peekedArg.Text
			if len(opt.Delimiter) == 0 {
				values = append(values, value)
			} else {
				values = append(values, strings.Split(value, opt.Delimiter)...)
			}
		}

		if opt.OverridePrev || params[opt.Key] == nil {
			params[opt.Key] = values
		} else {
			params[opt.Key] = append(params[opt.Key], values...)
		}
	}

	return &ParseResult{
		inputs:   args,
		params:   params,
		defaults: defaults,
		rests:    rests,
	}
}
