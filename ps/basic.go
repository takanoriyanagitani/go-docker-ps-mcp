package ps

import (
	"regexp"

	"github.com/docker/docker/api/types/container"
)

// BasicSummary represents the essential columns of "docker ps"
type BasicSummary struct {
	// Full container ID
	ID string `json:"Id"`

	// List of container names
	Names []string `json:"Names"`

	// Name of the image used
	Image string `json:"Image"`

	// Command string being run
	Command string `json:"Command"`

	// Unix timestamp
	Created int64 `json:"Created"`

	// Human-readable status (e.g., "Up 2 minutes")
	Status string `json:"Status"`

	// Machine-readable state (e.g., "running", "exited")
	State string `json:"State"`

	// List of port mappings
	Ports []container.Port `json:"Ports"`
}

type ContainerSummary container.Summary

func (s ContainerSummary) ToBasic() BasicSummary {
	var c container.Summary = container.Summary(s)
	return BasicSummary{
		ID:      c.ID,
		Names:   c.Names,
		Image:   c.Image,
		Command: c.Command,
		Created: c.Created,
		Status:  c.Status,
		State:   c.State,
		Ports:   c.Ports,
	}
}

type SummaryList []container.Summary

func (l SummaryList) ToBasicList(filter SummaryFilter) []BasicSummary {
	var src []container.Summary = l
	if src == nil {
		return nil
	}

	dst := make([]BasicSummary, 0, len(src))

	for _, s := range src {
		var keep bool = filter(s)
		if !keep {
			continue
		}
		dst = append(dst, ContainerSummary(s).ToBasic())
	}
	return dst
}

type SummaryFilter func(container.Summary) (keep bool)

func SummaryFilterStatic(keep bool) SummaryFilter {
	return func(_ container.Summary) bool { return keep }
}

var SummaryFilterDefault SummaryFilter = SummaryFilterStatic(true)

type ContainerNamesFilter func(names []string) (keep bool)

func (f ContainerNamesFilter) ToSummaryFilter() SummaryFilter {
	return func(s container.Summary) (keep bool) {
		var names []string = s.Names
		return f(names)
	}
}

type ContainerNameFilter func(name string) (keep bool)

func (f ContainerNameFilter) ToContainerNamesFilter() ContainerNamesFilter {
	return func(names []string) (keep bool) {
		for _, name := range names {
			return f(name)
		}

		return false
	}
}

func (f ContainerNameFilter) ToSummaryFilter() SummaryFilter {
	return f.ToContainerNamesFilter().ToSummaryFilter()
}

type ContainerNamePattern struct{ *regexp.Regexp }

func (p ContainerNamePattern) ToNameFilter() ContainerNameFilter {
	return func(name string) (keep bool) {
		return p.Regexp.MatchString(name)
	}
}

func (p ContainerNamePattern) ToSummaryFilter() SummaryFilter {
	return p.
		ToNameFilter().
		ToSummaryFilter()
}

type ContainerNamePatternString string

func (s ContainerNamePatternString) ToPattern() (ContainerNamePattern, error) {
	var expr string = string(s)
	pat, e := regexp.Compile(expr)
	return ContainerNamePattern{Regexp: pat}, e
}

func (s ContainerNamePatternString) ToSummaryFilter() (SummaryFilter, error) {
	pat, e := s.ToPattern()
	if nil != e {
		return nil, e
	}

	return pat.ToSummaryFilter(), nil
}
