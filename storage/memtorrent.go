package storage

import (
	"io"
	"strconv"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"github.com/patrickmn/go-cache"

	"github.com/thank243/trafficConsume/common/fakefile"
)

func (m *memTorrent) Piece(p metainfo.Piece) storage.PieceImpl {
	key := p.Hash().String()
	if m.isFake {
		key = key + "_" + strconv.Itoa(p.Index())
	}
	// Check if Piece in cache regardless
	if pie, ok := PieceCache().Get(key); ok {
		return pie.(*memTorrent)
	}

	// If not in cache, create new (depending on flags)
	return m.createAndCachePiece(p, key)
}

func (m *memTorrent) createAndCachePiece(p metainfo.Piece, key string) storage.PieceImpl {
	expiration := cache.DefaultExpiration
	mp := &memTorrent{
		ih:     p.Hash(),
		isFake: m.isFake,
	}

	if m.isFake {
		expiration = cache.NoExpiration
		mp.isCompleted = fakefile.R.Float64() < 0.8
	}

	PieceCache().Set(key, mp, expiration)
	return mp
}

func (m *memTorrent) Close() error {
	PieceCache().Delete(m.ih.String())
	return nil
}

func (m *memTorrent) ReadAt(p []byte, off int64) (n int, err error) {
	// fake file for upload
	if m.isFake {
		if len(p) != 0 {
			for i := range p {
				p[i] = 0xff
			}
			return len(p), nil
		}
	}

	return 0, io.EOF
}

func (m *memTorrent) WriteAt(p []byte, off int64) (n int, err error) {
	return len(p), nil
}

func (m *memTorrent) MarkComplete() error {
	return nil
}

func (m *memTorrent) MarkNotComplete() error {
	return nil
}

func (m *memTorrent) Completion() storage.Completion {
	if m.isFake {
		return storage.Completion{
			Complete: m.isCompleted,
			Ok:       true,
		}
	}

	return storage.Completion{
		Complete: false,
		Ok:       true,
	}
}

func (m *memTorrent) SelfHash() (metainfo.Hash, error) {
	return m.ih, nil
}
