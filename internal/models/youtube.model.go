package models

type YouTubeResponse struct {
	Items []YouTubeItem `json:"items"`
}

type YouTubeItem struct {
	Snippet YouTubeSnippet `json:"snippet"`
	ContentDetails YouTubeContentDetails `json:"contentDetails"`
}

type YouTubeSnippet struct {
	Title string `json:"title"`
	Description string `json:"description"`

	ChannelTitle string `json:"channelTitle"`

	PublishedAt string `json:"publishedAt"`

	Thumbnails YouTubeThumbnails `json:"thumbnails"`
}

type YouTubeThumbnails struct {
	High YouTubeThumbnail `json:"high"`
}

type YouTubeThumbnail struct {
	URL string `json:"url"`
}

type YouTubeContentDetails struct {
	Duration string `json:"duration"`
}