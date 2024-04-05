package main

import (
	"sync/atomic"

	aslog "github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	"github.com/thank243/trafficConsume/app/client"
	"github.com/thank243/trafficConsume/common/metahash"
	"github.com/thank243/trafficConsume/storage"
)

func main() {
	getVersion()

	cfg := torrent.NewDefaultClientConfig()
	cfg.DefaultStorage = new(storage.Client)
	cfg.ExtendedHandshakeClientVersion = string([]byte{0xde, 0xad, 0xbe, 0xef})
	cfg.Logger = aslog.Default.WithFilterLevel(aslog.Error)

	cli, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// Monitor stats
	go cli.Monitor()

	// add fake file
	cli.AddFakeTorrent()

	// // default torrents
	cli.AddTorrents(metahash.GetDefaultMetaHashes())

	// nyaa torrents
	nyaaHs := metahash.GetNyaaMetaHashes()
	cli.AddTorrents(nyaaHs)

	// cron jobs
	var isRunning atomic.Bool
	job := cron.New()
	_, err = job.AddFunc("@every 1h", func() {
		if isRunning.Load() {
			return
		}
		isRunning.Store(true)
		defer isRunning.Store(false)

		newNyaaHs := metahash.GetNyaaMetaHashes()
		cli.AddTorrents(newNyaaHs)

		dropHs := metahash.NeedDropTorrents(nyaaHs, newNyaaHs)
		for _, h := range dropHs {
			if t, ok := cli.Torrent(h); ok {
				t.Drop()
			}
		}
		log.Infof("Add %d torrents, Drop %d torrents", len(newNyaaHs)+len(dropHs)-len(nyaaHs), len(dropHs))
		nyaaHs = newNyaaHs
	})
	if err != nil {
		log.Fatal(err)
	}
	job.Start()

	// block wait completed
	cli.WaitAll()
}
