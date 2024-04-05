package client

import (
	"sync"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	log "github.com/sirupsen/logrus"

	"github.com/thank243/trafficConsume/common/fakefile"
	"github.com/thank243/trafficConsume/infra"
	"github.com/thank243/trafficConsume/storage"
)

func New(cfg *torrent.ClientConfig) (*Client, error) {
	cl, err := torrent.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	now := time.Now()

	return &Client{
		Client:            cl,
		totalStats:        stats{createdAt: now},
		fakeUploadStats:   stats{createdAt: now},
		fakeDownloadStats: stats{createdAt: now},
	}, nil
}

func (c *Client) AddTorrents(mhs []metainfo.Hash) {
	// default tracker servers
	trs := []string{
		"http://nyaa.tracker.wf:7777/announce",
		"http://p4p.arenabg.com:1337/announce",
		"udp://tracker.opentrackr.org:1337/announce",
	}

	var wg sync.WaitGroup
	for i := range mhs {
		wg.Add(1)
		go func(i int) {
			t, _ := c.AddTorrentInfoHash(mhs[i])
			t.AddTrackers([][]string{trs})
			if t.Info() == nil {
				<-t.GotInfo()
			}
			t.DownloadAll()
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func (c *Client) Monitor() {
	for now := range time.Tick(time.Second * 10) {
		totalBytes := c.ConnStats().BytesRead
		totalSpeed := c.speed(&c.totalStats, totalBytes, now)
		fakeDownSpeed, fakeUpSpeed, actPeers := c.torrentStats(now)

		log.Infof("Throughput: %s, Total: ↓ %s/s, Private: ↑ %s/s - ↓ %s/s, Pieces: %d, Peers: %d, Tasks: %d",
			infra.ByteCountIEC(totalBytes.Int64()),
			infra.ByteCountIEC(totalSpeed), infra.ByteCountIEC(fakeUpSpeed), infra.ByteCountIEC(fakeDownSpeed),
			storage.PieceCache().ItemCount(), actPeers, len(c.Torrents()))
	}
}

func (c *Client) torrentStats(now time.Time) (fakeDownSpeed int64, fakeUpSpeed int64, actPeers int) {
	for _, t := range c.Torrents() {
		actPeers += t.Stats().ActivePeers

		if t.InfoHash().String() == storage.FakeFileHash {
			fakeDownSpeed = c.speed(&c.fakeDownloadStats, t.Stats().BytesRead, now)
			fakeUpSpeed = c.speed(&c.fakeUploadStats, t.Stats().BytesWritten, now)
		}
	}
	return
}

func (c *Client) speed(s *stats, nowBytes torrent.Count, now time.Time) int64 {
	b := nowBytes.Int64()
	speed := (b - s.bytesCount) * 1000 / now.Sub(s.createdAt).Milliseconds()
	s.bytesCount = b
	s.createdAt = now
	return speed
}

func (c *Client) AddFakeTorrent() {
	f := &fakefile.FakeFile{
		Size:     1<<30 + 114514,
		FillByte: 0xff,
	}

	t, _ := c.AddTorrent(&metainfo.MetaInfo{
		InfoBytes: bencode.MustMarshal(f.BuildFakeFileInfo()),
	})

	trs := []string{
		"http://p4p.arenabg.com:1337/announce",
		"udp://tracker.opentrackr.org:1337/announce",
	}
	t.AddTrackers([][]string{trs})

	t.DownloadAll()
}
