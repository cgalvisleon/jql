package jql

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type From struct {
	Database string `json:"database"`
	Schema   string `json:"schema"`
	Name     string `json:"name"`
	As       string `json:"as"`
}

type Field struct {
	TypeColumn TypeColumn  `json:"type_column"`
	From       From        `json:"from"`
	Name       interface{} `json:"name"`
	As         string      `json:"as"`
}

func (s *Field) AS() string {
	return fmt.Sprintf("%s.%s", s.From.As, s.As)
}

/**
* FindField
* @param froms []*From, name string // from.name:as|1:30
* @return *Field
**/
func FindField(froms []*Froms, name string) *Field {
	pattern1 := regexp.MustCompile(`^([A-Za-z0-9]+)\.([A-Za-z0-9]+):([A-Za-z0-9]+)$`) // from.name:as
	pattern2 := regexp.MustCompile(`^([A-Za-z0-9]+)\.([A-Za-z0-9]+)$`)                // from.name
	pattern3 := regexp.MustCompile(`^([A-Za-z]+)\((.+)\):([A-Za-z0-9]+)$`)            // args(field):as
	pattern4 := regexp.MustCompile(`^([A-Za-z]+)\((.+)\)`)                            // args(field)
	pattern5 := regexp.MustCompile(`^(\d+)\|(\d+)$`)                                  // page:rows

	split := strings.Split(name, "|")
	if len(split) == 2 {
		name = split[0]
		limit := split[1]
		result := FindField(froms, name)
		if result != nil {
			if pattern5.MatchString(limit) {
				matches := pattern5.FindStringSubmatch(limit)
				if len(matches) == 3 {
					page, err := strconv.Atoi(matches[1])
					if err != nil {
						page = 0
					}
					rows, err := strconv.Atoi(matches[2])
					if err != nil {
						rows = 0
					}
					result.Page = page
					result.Rows = rows
				}
			}
		}

		return result
	}

	if pattern1.MatchString(name) {
		matches := pattern1.FindStringSubmatch(name)
		if len(matches) == 4 {
			fromName := matches[1]
			name = matches[2]
			as := matches[3]
			var result *Field
			for _, from := range froms {
				if from.As == fromName {
					result = from.FindField(name)
				} else if from.Model.Name == fromName {
					result = from.FindField(name)
				}
				if result != nil {
					result.As = as
					return result
				}
			}
		}
	} else if pattern2.MatchString(name) {
		matches := pattern2.FindStringSubmatch(name)
		if len(matches) == 3 {
			fromName := matches[1]
			name = matches[2]
			as := matches[2]
			var result *Field
			for _, from := range froms {
				if from.As == fromName {
					result = from.FindField(name)
				} else if from.Model.Name == fromName {
					result = from.FindField(name)
				}
				if result != nil {
					result.As = as
					return result
				}
			}
		}
	} else if pattern3.MatchString(name) {
		matches := pattern3.FindStringSubmatch(name)
		if len(matches) == 4 {
			aggregation := matches[1]
			name = matches[2]
			as := matches[3]
			result := FindField(froms, name)
			if result != nil {
				result.TypeColumn = TpAggregation
				result.Aggregation = GetAggregation(aggregation)
				result.As = as
				return result
			}
		}
	} else if pattern4.MatchString(name) {
		matches := pattern4.FindStringSubmatch(name)
		if len(matches) == 3 {
			aggregation := matches[1]
			name = matches[2]
			result := FindField(froms, name)
			if result != nil {
				result.TypeColumn = TpAggregation
				result.Aggregation = GetAggregation(aggregation)
				result.As = aggregation
				return result
			}
		}
	} else {
		if len(froms) == 0 {
			return nil
		}
		from := froms[0]
		result := from.FindField(name)
		if result != nil {
			return result
		}
	}

	return nil
}
