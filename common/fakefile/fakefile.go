package fakefile

import (
	"crypto/sha1"
	"io"

	"github.com/anacrolix/torrent/metainfo"
)

type FakeFile struct {
	Size     int
	FillByte byte
}

func (f *FakeFile) BuildFakePieces(pieceLen int64) []byte {
	var cumB []byte
	cycleNum := int64(f.Size) / pieceLen
	r := int64(f.Size) % pieceLen
	p := make([]byte, pieceLen)
	for i := range p {
		p[i] = f.FillByte
	}
	b := sha1.Sum(p)
	rb := sha1.Sum(p[:r])
	for i := int64(0); i < cycleNum; i++ {
		cumB = append(cumB, b[:]...)
	}
	cumB = append(cumB, rb[:]...)

	return cumB
}

func (f *FakeFile) Read(p []byte) (n int, err error) {
	readSize := len(p)
	r := f.Size
	if r < readSize {
		readSize = r
	}

	for i := 0; i < readSize; i++ {
		p[i] = f.FillByte
	}

	r -= readSize
	if r <= 0 {
		err = io.EOF
	}

	return readSize, err
}

func (f *FakeFile) BuildFakeFileInfo() *metainfo.Info {
	info := &metainfo.Info{
		Name:        "fake.file",
		Length:      int64(f.Size),
		PieceLength: metainfo.ChoosePieceLength(int64(f.Size)),
	}
	info.Pieces = f.BuildFakePieces(info.PieceLength)

	return info
}
