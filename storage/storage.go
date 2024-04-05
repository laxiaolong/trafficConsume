package storage

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
)

func (c *Client) OpenTorrent(info *metainfo.Info, infoHash metainfo.Hash) (storage.TorrentImpl, error) {
	mp := &memTorrent{
		ih:     infoHash, // for whole file hash not piece hash
		isFake: infoHash.String() == FakeFileHash,
	}

	return storage.TorrentImpl{
		Piece: mp.Piece,
	}, nil
}
