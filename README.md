> [!NOTE]
> This project is a fork of [n10v/id3v2](https://github.com/n10v/id3v2), which has not been actively maintained.  
> Support for synced lyrics has been added, and fixes will be applied as needed.

# id3v2

Implementation of ID3 v2.3 and v2.4 in native Go.

## Installation

```
go get -u github.com/joseavilasg/go-id3v2
```

## Documentation

https://pkg.go.dev/github.com/joseavilasg/go-id3v2

## Usage example

```go
package main

import (
	"fmt"
	"log"

	"github.com/joseavilasg/go-id3v2"
)

func main() {
	tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
	if err != nil {
 		log.Fatal("Error while opening mp3 file: ", err)
 	}
	defer tag.Close()

	// Read tags
	fmt.Println(tag.Artist())
	fmt.Println(tag.Title())

	// Set tags
	tag.SetArtist("Aphex Twin")
	tag.SetTitle("Xtal")

	comment := id3v2.CommentFrame{
		Encoding:    id3v2.EncodingUTF8,
		Language:    "eng",
		Description: "My opinion",
		Text:        "I like this song!",
	}
	tag.AddCommentFrame(comment)

	// Write tag to file.mp3
	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}
}
```

## Read multiple frames

```go
pictures := tag.GetFrames(tag.CommonID("Attached picture"))
for _, f := range pictures {
	pic, ok := f.(id3v2.PictureFrame)
	if !ok {
		log.Fatal("Couldn't assert picture frame")
	}

	// Do something with picture frame
	fmt.Println(pic.Description)
}
```

## Encodings

For example, if you want to set comment frame with custom encoding,
you may do the following:

```go
comment := id3v2.CommentFrame{
	Encoding:    id3v2.EncodingUTF16,
	Language:    "ger",
	Description: "Tier",
	Text:        "Der Löwe",
}
tag.AddCommentFrame(comment)
```

`Text` field will be automatically encoded with UTF-16BE with BOM and written to w.

UTF-8 is default for v2.4, ISO-8859-1 - for v2.3.
