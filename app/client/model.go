package client

import (
	"time"

	"github.com/anacrolix/torrent"
)

type Client struct {
	*torrent.Client
	totalStats        stats
	fakeUploadStats   stats
	fakeDownloadStats stats
}

type stats struct {
	bytesCount int64
	createdAt  time.Time
}
