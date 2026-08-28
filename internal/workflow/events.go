package workflow

import (
	"sort"

	"example.com/graduation-showcase/internal/domain"
)

type Event struct {
	Step string
	At   int64
	Note string
}

type EventLog struct {
	WorkflowID string
	Events     []Event
}

func NewEventLog(workflowID string) EventLog {
	return EventLog{WorkflowID: workflowID, Events: []Event{}}
}

func (l *EventLog) Append(step, note string, at int64) {
	l.Events = append(l.Events, Event{Step: step, Note: note, At: at})
}

func (l EventLog) Steps() []string {
	result := make([]string, 0, len(l.Events))
	for _, event := range l.Events {
		result = append(result, event.Step)
	}
	return result
}

func (l EventLog) Last() Event {
	if len(l.Events) == 0 {
		return Event{}
	}
	return l.Events[len(l.Events)-1]
}

func (l EventLog) Ordered() EventLog {
	result := EventLog{WorkflowID: l.WorkflowID, Events: append([]Event(nil), l.Events...)}
	sort.SliceStable(result.Events, func(i, j int) bool {
		if result.Events[i].At != result.Events[j].At {
			return result.Events[i].At < result.Events[j].At
		}
		return result.Events[i].Step < result.Events[j].Step
	})
	return result
}

func (l EventLog) ToWorkflow(now int64, name string) domain.Workflow {
	ordered := l.Ordered()
	status := "running"
	if len(ordered.Events) >= 4 {
		status = "completed"
	}
	return domain.Workflow{ID: l.WorkflowID, Name: name, Status: status, Steps: ordered.Steps(), StartedAt: ordered.FirstAt(), CompletedAt: now}
}

func (l EventLog) FirstAt() int64 {
	if len(l.Events) == 0 {
		return 0
	}
	return l.Ordered().Events[0].At
}
