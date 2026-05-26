package query

import (
	"net/url"
	"regexp"
	"strings"
)

type Filter struct {
	Field    string
	Operator string
	Value    string
}

type SortField struct {
	Field string
	Desc  bool
}

type QueryParams struct {
	Filters    []Filter
	SortFields []SortField
	Page       int
	Limit      int
}

func ParseFilters(query url.Values) []Filter {
	var filters []Filter
	re := regexp.MustCompile(`^(\w+)\[(\w+)\]$`)

	for key, values := range query {
		if len(values) == 0 {
			continue
		}

		matches := re.FindStringSubmatch(key)
		if matches == nil {
			// Fallback: Handle exact match without operator e.g. status=active
			// But avoid parsing standard query params like page, limit, sort
			if isReservedParam(key) {
				continue
			}
			filters = append(filters, Filter{
				Field:    key,
				Operator: "eq",
				Value:    values[0],
			})

			continue
		}

		filters = append(filters, Filter{
			Field:    matches[1],
			Operator: matches[2],
			Value:    values[0],
		})
	}

	return filters
}

func ParseSort(raw string) []SortField {
	var sorts []SortField
	for _, s := range strings.Split(raw, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if strings.HasPrefix(s, "-") {
			sorts = append(sorts, SortField{Field: s[1:], Desc: true})
		} else {
			sorts = append(sorts, SortField{Field: s, Desc: false})
		}
	}

	return sorts
}

func isReservedParam(key string) bool {
	reserved := map[string]bool{
		"page":     true,
		"limit":    true,
		"per_page": true,
		"sort":     true,
	}

	return reserved[key]
}
