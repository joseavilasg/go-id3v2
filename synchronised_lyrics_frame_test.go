package id3v2

import (
	"bytes"
	"testing"
)

func TestSynchronisedLyricsFrame(t *testing.T) {
	contentDescriptor := "Test descriptor"
	syncedTexts := []SyncedText{
		{Text: "First line", Timestamp: 1000},
		{Text: "Second line", Timestamp: 2000},
	}

	frame := SynchronisedLyricsFrame{
		Encoding:          EncodingUTF8,
		Language:          "eng",
		TimestampFormat:   SYLTAbsoluteMillisecondsTimestampFormat,
		ContentType:       SYLTLyricsContentType,
		ContentDescriptor: contentDescriptor,
		SynchronizedTexts: syncedTexts,
	}

	buf := new(bytes.Buffer)

	if _, err := frame.WriteTo(buf); err != nil {
		t.Fatalf("Failed to write frame: %v", err)
	}

	parsed, err := parseSynchronisedLyricsFrame(newBufReader(buf), 4)
	if err != nil {
		t.Fatalf("Failed to parse frame: %v", err)
	}

	parsedFrame, ok := parsed.(SynchronisedLyricsFrame)
	if !ok {
		t.Fatalf("Parsed frame is not of type SynchronisedLyricsFrame")
	}

	if parsedFrame.ContentDescriptor != contentDescriptor {
		t.Errorf("Expected content descriptor: %q, got: %q", contentDescriptor, parsedFrame.ContentDescriptor)
	}

	if len(parsedFrame.SynchronizedTexts) != len(syncedTexts) {
		t.Fatalf("Expected %d synchronized texts, got %d", len(syncedTexts), len(parsedFrame.SynchronizedTexts))
	}

	for i, text := range syncedTexts {
		if parsedFrame.SynchronizedTexts[i].Text != text.Text || parsedFrame.SynchronizedTexts[i].Timestamp != text.Timestamp {
			t.Errorf("Mismatch in synchronized text at index %d: expected %+v, got %+v", i, text, parsedFrame.SynchronizedTexts[i])
		}
	}
}
