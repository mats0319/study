package main

type event struct {
	receivers []string
	data      *eventSource
}

type eventSource struct {
	id    string
	event string
	data  string
}

func newEvent(receivers []string, id, eventName, data string) *event {
	return &event{
		receivers: receivers,
		data: &eventSource{
			id:    id,
			event: eventName,
			data:  data,
		},
	}
}

func (ed *eventSource) format() string {
	if ed == nil {
		return ""
	}

	res := ""
	if len(ed.id) > 0 {
		res += "id: " + ed.id + "\n"
	}
	if len(ed.event) > 0 {
		res += "event: " + ed.event + "\n"
	}

	res += "data: " + ed.data + "\n\n"

	return res
}
