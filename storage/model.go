package storage

import (
	"github.com/anacrolix/torrent/metainfo"
)

const (
	FakeFileHash = "d9511e3a9284a9e1833b258fd4c33247c3cacfe4"
)

type Client struct {
}

type memTorrent struct {
	ih          metainfo.Hash
	isFake      bool
	isCompleted bool
}
