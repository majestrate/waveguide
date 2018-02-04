package video

type Info struct {
	VideoCodec string
	AudioCodec string
	Width      int
	Height     int
}

type StreamInfo struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	AVC       string `json:"is_avc"`
	Channels  int    `json:"channels"`
}

type ProbeInfo struct {
	Streams []StreamInfo `json:"streams"`
}

func (i StreamInfo) IsAudio() bool {
	return i.CodecType == "audio"
}

func (i StreamInfo) IsVideo() bool {
	return i.CodecType == "video"
}

func (i StreamInfo) IsAVC() bool {
	return i.AVC == "true"
}

func (i Info) Matches(infos []StreamInfo) bool {
	var foundVideo, foundAudio bool
	for _, info := range infos {
		if info.IsAudio() && !foundAudio {
			foundAudio = i.AudioCodec == info.CodecName
		} else if info.IsVideo() && !foundVideo {
			foundVideo = i.VideoCodec == info.CodecName
		}
	}
	return foundVideo && foundAudio
}
