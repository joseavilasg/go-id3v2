// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"encoding/binary"
	"io"
)

// SynchronisedLyricsFrame is used to work with SYLT frames.
// The information about how to add synchronised lyrics/text frame to tag
// you can see in the docs to tag.AddSynchronisedLyricsFrame function.
//
// You must choose a three-letter language code from
// ISO 639-2 code list: https://www.loc.gov/standards/iso639-2/php/code_list.php
type SynchronisedLyricsFrame struct {
	Encoding          Encoding
	Language          string
	TimestampFormat   byte
	ContentType       byte
	ContentDescriptor string
	SynchronizedTexts []SyncedText
}

const (
	SYLTAbsoluteMpegFramesTimestampFormat = iota + 1
	SYLTAbsoluteMillisecondsTimestampFormat
)

const (
	SYLTOtherContentType = iota
	SYLTLyricsContentType
	SYLTTextTranscriptionContentType
	SYLTMovementContentType
	SYLTEventsContentType
	SYLTChordContentType
	SYLTTriviaContentType
	SYLTWebpageURLsContentType
	SYLTImageURLsContentType
)

type SyncedText struct {
	Text      string
	Timestamp uint32
}

func (sylf SynchronisedLyricsFrame) Size() int {
	var s int
	for _, v := range sylf.SynchronizedTexts {
		s += encodedSize(v.Text, sylf.Encoding)
		s += len(sylf.Encoding.TerminationBytes)
		s += 4
	}

	return 1 + len(sylf.Language) + encodedSize(sylf.ContentDescriptor, sylf.Encoding) +
		+len(sylf.Encoding.TerminationBytes) + s +
		+1 + 1
}

func (sylf SynchronisedLyricsFrame) UniqueIdentifier() string {
	return sylf.Language + sylf.ContentDescriptor
}

func (sy SyncedText) uintToBytes() []byte {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, sy.Timestamp)
	return bs
}

func bytesToInt(timeStampBytes []byte) uint32 {
	return binary.BigEndian.Uint32(timeStampBytes)
}

func (sylf SynchronisedLyricsFrame) WriteTo(w io.Writer) (n int64, err error) {
	if len(sylf.Language) != 3 {
		return n, ErrInvalidLanguageLength
	}
	return useBufWriter(w, func(bw *bufWriter) {
		bw.WriteByte(sylf.Encoding.Key)
		bw.WriteString(sylf.Language)
		bw.WriteByte(sylf.TimestampFormat)
		bw.WriteByte(sylf.ContentType)
		bw.EncodeAndWriteText(sylf.ContentDescriptor, sylf.Encoding)
		bw.Write(sylf.Encoding.TerminationBytes)
		for _, v := range sylf.SynchronizedTexts {
			bw.EncodeAndWriteText(v.Text, sylf.Encoding)
			bw.Write(sylf.Encoding.TerminationBytes)
			bw.Write(v.uintToBytes())
		}
	})
}

func parseSynchronisedLyricsFrame(br *bufReader, version byte) (Framer, error) {
	encoding := getEncoding(br.ReadByte())
	language := br.Next(3)
	timestampFormat := br.ReadByte()
	contentType := br.ReadByte()
	contentDescriptor := br.ReadText(encoding)

	if br.Err() != nil {
		return nil, br.Err()
	}

	var s []SyncedText
	for {
		textLyric, err := br.readTillDelims(encoding.TerminationBytes)
		if err != nil {
			break
		}
		t := SyncedText{Text: decodeText(textLyric, encoding)}
		br.Next(len(encoding.TerminationBytes))
		timeStamp := br.Next(4)
		timeStampUint := bytesToInt(timeStamp)
		t.Timestamp = timeStampUint
		s = append(s, t)
	}
	sylf := SynchronisedLyricsFrame{
		Encoding:          encoding,
		Language:          string(language),
		TimestampFormat:   timestampFormat,
		ContentType:       contentType,
		ContentDescriptor: decodeText(contentDescriptor, encoding),
		SynchronizedTexts: s,
	}

	return sylf, nil
}
