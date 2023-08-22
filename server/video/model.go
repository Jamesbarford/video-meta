package video

/* Essentially what AWS does for S3 Tagging */
type VideoMetaDataPayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type VideoMetaData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Id    int    `json:"id"`
}
