package identifier

import (
	"enforce-tool-versions/command"
	"errors"
	"github.com/rs/zerolog"
	"strings"
)

// enum for Programs that we can identify
type Program int

const (
	Make Program = iota
	Git
)

// Version is a string, should be lexicographically sortable (e.g. semver).
type Version string

var identifierMap = map[Program]func(string, *zerolog.Logger) (Version, error){
	Make: identifyMake,
	Git:  identifyGit,
}

var programMap = map[Program]string{
	Make: "make",
	Git:  "git",
}

// GetProgramName returns the name of the program p.
func GetProgramName(p Program) string {
	return programMap[p]
}

var (
	ErrProgramNotSupported = errors.New("program not supported")
)

// Identify returns the version of the program p, or an error if the program is not supported.
func Identify(p Program, zlog *zerolog.Logger) (Version, error) {
	identifier, ok := identifierMap[p]
	if !ok {
		zlog.Debug().Msg("program not supported")
		return "", ErrProgramNotSupported
	}
	versionOutput, err := getProgramVersionOutput(p, zlog)
	if err != nil {
		zlog.Debug().Err(err).Msg("failed to get program version output")
		return "", err
	}
	return identifier(versionOutput, zlog)
}

// s is a single line, e.g.
//
// git version 2.39.1
func identifyGit(s string, zlog *zerolog.Logger) (Version, error) {
	word, err := getLastWordOnFirstLine(s)
	if err != nil {
		zlog.Debug().Err(err).Msg("failed to get last word on first line")
		return "", err
	}
	return Version(word), nil
}

// On the first line, get the last whitespace-delimited element.
//
// Example s:
//
// GNU Make 4.4
// Built for aarch64-apple-darwin21.6.0
// Copyright (C) 1988-2022 Free Software Foundation, Inc.
// License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>
// This is free software: you are free to change and redistribute it.
// There is NO WARRANTY, to the extent permitted by law.
func identifyMake(s string, zlog *zerolog.Logger) (Version, error) {
	word, err := getLastWordOnFirstLine(s)
	if err != nil {
		zlog.Debug().Err(err).Msg("failed to get last word on first line")
		return "", err
	}
	return Version(word), nil
}

func getLastWordOnFirstLine(s string) (string, error) {
	lines := strings.Split(s, "\n")
	if len(lines) == 0 {
		return "", errors.New("no lines in output")
	}
	words := strings.Split(lines[0], " ")
	if len(words) == 0 {
		return "", errors.New("no words in first line")
	}
	return words[len(words)-1], nil
}

func getProgramVersionOutput(p Program, zlog *zerolog.Logger) (string, error) {
	output, err := command.RunCommand(GetProgramName(p), "--version")
	if err != nil {
		zlog.Debug().Str("output", output).Err(err).Msg("failed to run command")
		return "", err
	}
	return output, nil
}
