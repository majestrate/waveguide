package streaming

const MaxMagnets = 5

type StreamInfo struct {
	Magnets []string
}

func (i *StreamInfo) Add(url string) {
	if len(i.Magnets) > MaxMagnets {
		i.Magnets = append(i.Magnets[1:], url)
	} else {
		i.Magnets = append(i.Magnets, url)
	}
}
