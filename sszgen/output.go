package sszgen

import "io"

type MergedFile struct {
	sats []SSZSatisfier
}

func (mf *MergedFile) Accumulate(s SSZSatisfier) {
	mf.sats = append(mf.sats, s)
}

func (oa *MergedFile) Merge() (io.Reader, error) {
	return nil, nil
}
