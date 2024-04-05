package metahash

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/types/infohash"
	"github.com/imroc/req/v3"
	log "github.com/sirupsen/logrus"
)

func GetNyaaMetaHashes() []metainfo.Hash {
	client := req.C()
	resp, err := client.R().Get("https://sukebei.nyaa.si/?s=seeders&o=desc")
	if err != nil {
		log.Error(err)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Error(err)
		return nil
	}

	var els []string
	doc.Find(`a[href^="magnet"]`).Each(func(i int, s *goquery.Selection) {
		if val, ok := s.Attr("href"); ok {
			u, err := url.Parse(val)
			if err != nil {
				log.Error(err)
				return
			}

			q := u.Query()
			els = append(els, strings.Split(q.Get("xt"), ":")[2])
		}
	})

	length := 5
	if len(els) < length {
		length = len(els)
	}
	log.Infof("Get %d meta info hash, add the top %d to task", len(els), length)

	hs := make([]metainfo.Hash, length)
	for i := 0; i < length; i++ {
		hs[i] = infohash.FromHexString(els[i])
	}

	return hs
}

// GetDefaultMetaHashes default torrents
func GetDefaultMetaHashes() []metainfo.Hash {
	return []metainfo.Hash{
		infohash.FromHexString("76c23875f5d146993d79170e78e2de43c3178cf1"),
		infohash.FromHexString("5a29a7e432691d8f590c1b8e05b6678eee72fe68"),
		infohash.FromHexString("6dc8bc544faa0216c044855ddc51cd9252da6a2b"),
	}
}

func NeedDropTorrents(old []metainfo.Hash, new []metainfo.Hash) []metainfo.Hash {
	if len(new) == 0 {
		return nil
	}

	newMap := make(map[metainfo.Hash]bool)
	for i := range new {
		newMap[new[i]] = true
	}

	var dropTorrents []metainfo.Hash
	for i := range old {
		if _, ok := newMap[old[i]]; !ok {
			dropTorrents = append(dropTorrents, old[i])
		}
	}
	return dropTorrents
}
