package upstream

import (
	"regexp"
	"strings"
)

type QueryParser func(query string) []Resource

var (
	tableRefWithWildcard = "(?:`([\\w-]+)`\\.([\\w-]+)\\.([\\w-\\*?]+)|`?([\\w-]+)\\.([\\w-]+)\\.([\\w-\\*?]+)`?)"
	tableRef             = "(?:`([\\w-]+)`\\.([\\w-]+)\\.([\\w-]+)|`?([\\w-]+)\\.([\\w-]+)\\.([\\w-]+)`?)"

	topLevelUpstreamsPattern = regexp.MustCompile("" +
		"(?i)(?:FROM)\\s*(?:/\\*\\s*([a-zA-Z0-9@_-]*)\\s*\\*/)?\\s+" + tableRefWithWildcard +
		"|" +
		"(?i)(?:JOIN)\\s*(?:/\\*\\s*([a-zA-Z0-9@_-]*)\\s*\\*/)?\\s+" + tableRef +
		"|" +
		"(?i)(?:WITH)\\s*(?:/\\*\\s*([a-zA-Z0-9@_-]*)\\s*\\*/)?\\s+" + tableRef + "\\s+(?:AS)" +
		"|" +
		"(?i)(?:VIEW)\\s*(?:/\\*\\s*([a-zA-Z0-9@_-]*)\\s*\\*/)?\\s+" + tableRef +
		"|" +
		"(?i)(?:/\\*\\s*([a-zA-Z0-9@_-]*)\\s*\\*/)?\\s+`([\\w-]+)\\.([\\w-]+)\\.([\\w-]+)`\\s*(?:AS)?")

	singleLineCommentsPattern = regexp.MustCompile(`(--.*)`)
	multiLineCommentsPattern  = regexp.MustCompile(`(((/\*)+?[\w\W]*?(\*/)+))`)
	specialCommentPattern     = regexp.MustCompile(`(/\*\s*(@[a-zA-Z0-9_-]+)\s*\*/)`)
)

func ParseTopLevelUpstreamsFromQuery(query string) []Resource {
	cleanedQuery := cleanQueryFromComment(query)

	resourcesFound := make(map[Resource]bool)
	pseudoResources := make(map[Resource]bool)

	matches := topLevelUpstreamsPattern.FindAllStringSubmatch(cleanedQuery, -1)
	for _, match := range matches {
		var ignoreUpstreamIdx int
		var projectIdx, datasetIdx, nameIdx []int
		tokens := strings.Fields(match[0])
		clause := strings.ToLower(tokens[0])

		switch clause {
		case "from":
			ignoreUpstreamIdx = 1
			projectIdx, datasetIdx, nameIdx = []int{2, 5}, []int{3, 6}, []int{4, 7}
		case "join":
			ignoreUpstreamIdx = 8
			projectIdx, datasetIdx, nameIdx = []int{9, 12}, []int{10, 13}, []int{11, 14}
		case "with":
			ignoreUpstreamIdx = 15
			projectIdx, datasetIdx, nameIdx = []int{16, 19}, []int{17, 20}, []int{18, 21}
		case "view":
			ignoreUpstreamIdx = 22
			projectIdx, datasetIdx, nameIdx = []int{23, 26}, []int{24, 27}, []int{25, 28}
		default:
			ignoreUpstreamIdx = 29
			projectIdx, datasetIdx, nameIdx = []int{30}, []int{31}, []int{32}
		}

		project := firstNonEmptySubmatch(match, projectIdx)
		dataset := firstNonEmptySubmatch(match, datasetIdx)
		name := firstNonEmptySubmatch(match, nameIdx)

		if project == "" || dataset == "" || name == "" {
			continue
		}

		if strings.TrimSpace(match[ignoreUpstreamIdx]) == "@ignoreupstream" {
			continue
		}

		if clause == "view" {
			continue
		}

		resource := Resource{
			Project: project,
			Dataset: dataset,
			Name:    name,
		}

		if clause == "with" {
			pseudoResources[resource] = true
		} else {
			resourcesFound[resource] = true
		}
	}

	var output []Resource

	for resource := range resourcesFound {
		if pseudoResources[resource] {
			continue
		}
		output = append(output, resource)
	}

	return output
}

func cleanQueryFromComment(query string) string {
	cleanedQuery := singleLineCommentsPattern.ReplaceAllString(query, "")

	matches := multiLineCommentsPattern.FindAllString(query, -1)
	for _, match := range matches {
		if specialCommentPattern.MatchString(match) {
			continue
		}
		cleanedQuery = strings.ReplaceAll(cleanedQuery, match, "")
	}

	return cleanedQuery
}

func firstNonEmptySubmatch(match []string, indices []int) string {
	for _, idx := range indices {
		if idx < len(match) && match[idx] != "" {
			return match[idx]
		}
	}
	return ""
}
