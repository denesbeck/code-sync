package main

import (
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// default: all files are allowed (tracked).
// ignore: defines patterns to exclude.
// allow: defines patterns to re-include (override an ignore).

type Rules struct {
	Ignore []string `yaml:"ignore"`
	Allow  []string `yaml:"allow"`
}

func readRules() (*Rules, error) {
	Debug("Reading rules...")
	_, err := os.Stat(".nexio.rules.yml")
	if os.IsNotExist(err) {
		Debug("\".nexio.rules.yml\" doesn't exist.")
		return nil, err
	}

	if err != nil {
		Debug("Unable to get metadata for \".nexio.rules.yml\".")
		MustSucceed(err, "operation failed")
		return nil, err
	}

	var content Rules
	rulesFile, err := os.ReadFile(".nexio.rules.yml")
	if err != nil {
		Debug("Unable to read \".nexio.rules.yml\" file.")
		MustSucceed(err, "operation failed")
		return nil, err
	}
	if err = yaml.Unmarshal(rulesFile, &content); err != nil {
		Debug("Unable to unmarshal `.nexio.rules.yml` file.")
		MustSucceed(err, "operation failed")
		return nil, err
	}
	Debug("`.nexio.rules.yml` found.")
	Debug("Content: %+v", content)
	return &content, nil
}

func pathToRegexp() (ignore []*regexp.Regexp, allow []*regexp.Regexp, err error) {
	rules, err := readRules()
	if err != nil {
		return nil, nil, err
	}

	var ignoreRegexps []*regexp.Regexp
	for _, rule := range rules.Ignore {
		pattern, err := regexp.Compile(rule)
		if err != nil {
			patternRegexp, err := patternToRegexp(rule)
			if err != nil {
				return nil, nil, err
			}
			ignoreRegexps = append(ignoreRegexps, patternRegexp)
		} else {
			ignoreRegexps = append(ignoreRegexps, pattern)
		}
	}

	var allowRegexps []*regexp.Regexp
	for _, rule := range rules.Allow {
		pattern, err := regexp.Compile(rule)
		if err != nil {
			patternRegexp, err := patternToRegexp(rule)
			if err != nil {
				return nil, nil, err
			}
			allowRegexps = append(allowRegexps, patternRegexp)
		} else {
			allowRegexps = append(allowRegexps, pattern)
		}
	}

	return ignoreRegexps, allowRegexps, nil
}

func patternToRegexp(pattern string) (*regexp.Regexp, error) {
	regexpPattern := strings.ReplaceAll(pattern, "**", "__DOUBLE_STAR__")
	regexpPattern = strings.ReplaceAll(regexpPattern, "*", "__STAR__")
	regexpPattern = regexp.QuoteMeta(regexpPattern)
	regexpPattern = strings.ReplaceAll(regexpPattern, "__STAR__", "[^/]*")
	regexpPattern = strings.ReplaceAll(regexpPattern, "__DOUBLE_STAR__", ".*")

	isPath := strings.Contains(pattern, "/")
	if isPath {
		regexpPattern = "^" + regexpPattern + "$"
	} else {
		regexpPattern = "(^|/)" + regexpPattern + "(/|$)"
	}
	return regexp.Compile(regexpPattern)
}

func ShouldIgnore(path string) bool {
	Debug("Checking if %s should be ignored...", path)
	ignoreRegexps, allowRegexps, err := pathToRegexp()
	if err != nil {
		if os.IsNotExist(err) {
			Debug("Rules file not found, not ignoring path: %s", path)
			return false
		}
		Debug("Error reading rules file: %v", err)
		return false
	}

	if len(allowRegexps) == 0 && len(ignoreRegexps) == 0 {
		return false
	}

	for _, pattern := range allowRegexps {
		if pattern.MatchString(path) {
			Debug("Path should not be ignored: %s", path)
			return false
		}
	}

	for _, pattern := range ignoreRegexps {
		if pattern.MatchString(path) {
			Debug("Path should be ignored: %s", path)
			return true
		}
	}
	Debug("Path should not be ignored: %s", path)
	return false
}
