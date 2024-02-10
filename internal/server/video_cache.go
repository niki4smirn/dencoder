package server

import (
	"bytes"
	"errors"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

type batches [][]byte

type VideoProvider struct {
	bs batches
}

type VideosCache struct {
	cache_data *expirable.LRU[string, VideoProvider]
}

var fileNotFound = errors.New("file not found")

func intoBatches(content []byte) batches {
	bs := make(batches, 0)
	for pos := 0; pos < len(content); pos += batchSize {
		r := min(pos+batchSize, len(content))
		bs = append(bs, content[pos:r])
	}
	return bs
}

func (v *VideosCache) Write(filename string, content []byte) {
	v.cache_data.Add(filename, VideoProvider{bs: intoBatches(content)})
}

func (v *VideosCache) Contains(filename string) bool {
	return v.cache_data.Contains(filename)
}

func (v *VideosCache) GetProvdier(filename string) (*VideoProvider, error) {
	vProvider, found := v.cache_data.Get(filename)
	if !found {
		return nil, fileNotFound
	}
	return &vProvider, nil
}

func countBytes(bs batches) uint64 {
	if len(bs) == 1 {
		return uint64(len(bs[0]))
	}
	return uint64((len(bs)-1)*batchSize + len(bs[len(bs)-1]))
}

var indexOutOfBounds = errors.New("index out of bounds")

// [from, to]
func (v *VideoProvider) Read(from, to uint64) ([]byte, error) {
	if to >= countBytes(v.bs) {
		return nil, indexOutOfBounds
	}
	fromBatch := from / batchSize
	toBatch := to/batchSize + 1
	content := bytes.Join(v.bs[fromBatch:toBatch], []byte{})
	partSize := to - from + 1
	from = from % batchSize
	content = content[from:(from + partSize)]
	return content, nil
}

func (v *VideoProvider) Size() uint64 {
	return countBytes(v.bs)
}
