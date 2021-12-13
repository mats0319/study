package main

type event struct {
    to   []string
    data *eventData
}

type eventData struct {
    id    string
    event string
    data  string
}

func newEvent(to []string, data string) *event {
    return &event{
        to: to,
        data: &eventData{
            data: data,
        },
    }
}

func (ed *eventData) format() string {
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
