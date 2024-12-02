package model

type Preference struct {
	Low    string `json:"low"`
	Medium string `json:"medium"`
	High   string `json:"high"`
}

// get channel
func (p *Preference) GetChannel(priority string) string {
	switch priority {
	case "low":
		return p.Low
	case "medium":
		return p.Medium
	case "high":
		return p.High
	default:
		return ""
	}
}
