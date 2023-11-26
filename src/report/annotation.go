package report

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/cultureamp/ecrscanresults/finding"
	"github.com/cultureamp/ecrscanresults/findingconfig"
	"github.com/cultureamp/ecrscanresults/registry"
	"github.com/justincampbell/timeago"
	"golang.org/x/exp/maps"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed annotation.gohtml
var annotationTemplateSource string

type AnnotationContext struct {
	Image                     registry.RegistryInfo
	ImageLabel                string
	FindingSummary            finding.Summary
	CriticalSeverityThreshold int32
	HighSeverityThreshold     int32
}

func (c AnnotationContext) Render() ([]byte, error) {
	t, err := template.
		New("annotation").
		Funcs(template.FuncMap{
			"hasUntilValue": hasUntilValue,
			"titleCase": func(s string) string {
				c := cases.Title(language.English)
				return c.String(s)
			},
			"lowerCase": strings.ToLower,
			"nbsp": func(input string) any {
				if len(input) > 0 {
					return input
				} else {
					return template.HTML(`&nbsp;`)
				}
			},
			"timeAgo": func(tm *time.Time) string {
				if tm == nil {
					return ""
				}

				return timeago.FromTime(*tm)
			},
			"sortFindings":   sortFindings,
			"sortSeverities": sortSeverities,
			"string": func(input any) (string, error) {
				if strg, ok := input.(fmt.Stringer); ok {
					return strg.String(), nil
				}

				return fmt.Sprintf("%s", input), nil
			},
		}).
		Parse(annotationTemplateSource)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, c)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// sortFindings sorts by severity rank, then CVSS3/CVSS2 score _descending_, then by CVE descending
func sortFindings(findings []finding.Detail) []finding.Detail {
	// shallow clone, don't affect source array
	sorted := slices.Clone(findings)

	slices.SortFunc(sorted, func(a, b finding.Detail) int {
		// first by severity rank
		sevRank := compareSeverities(a.Severity, b.Severity)
		if sevRank != 0 {
			return sevRank
		}

		// then by CVSS3 score
		cvss3 := compareCVSSScore(a.CVSS3, b.CVSS3)
		if cvss3 != 0 {
			return cvss3 * -1 // descending
		}

		// then by CVSS2 score
		cvss2 := compareCVSSScore(a.CVSS2, b.CVSS2)
		if cvss2 != 0 {
			return cvss2 * -1 // descending
		}

		// descending order of CVE, in general this means that newer CVEs will be at
		// the top
		return strings.Compare(b.Name, a.Name)
	})

	return sorted
}

func sortSeverities(severityCounts map[types.FindingSeverity]finding.SeverityCount) []types.FindingSeverity {
	// severities are the map key in the incoming data structure
	severities := maps.Keys(severityCounts)

	slices.SortFunc(severities, compareSeverities)

	return severities
}

// sort severity strings by rank, then alphabetically
func compareSeverities(a, b types.FindingSeverity) int {
	rank := rankSeverity(a) - rankSeverity(b)

	if rank != 0 {
		return rank
	}

	// for unknown severities, sort alphabetically
	return strings.Compare(string(a), string(b))
}

func rankSeverity(f types.FindingSeverity) int {
	switch f {
	case types.FindingSeverityCritical:
		return 0
	case types.FindingSeverityHigh:
		return 1
	case types.FindingSeverityMedium:
		return 2
	case types.FindingSeverityLow:
		return 3
	case types.FindingSeverityInformational:
		return 4
	case types.FindingSeverityUndefined:
		return 5
	}

	return 100
}

func compareCVSSScore(a, b finding.CVSSScore) int {
	switch {
	case a.Score == nil && b.Score == nil:
		return 0
	case a.Score == nil:
		return -1
	case b.Score == nil:
		return 1
	default:
		return a.Score.Cmp(*b.Score)
	}
}

func hasUntilValue(until findingconfig.UntilTime) bool {
	return time.Time(until).After(time.Time(findingconfig.UntilTime{}))
}
