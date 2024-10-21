package main

import (
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3" //nolint
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileFrom, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	buf, er := io.ReadAll(fileFrom)
	if er != nil {
		log.Fatal("не прочитал")
	}

	if len(buf) < int(offset) {
		return ErrOffsetExceedsFileSize
	}

	errOp := fileFrom.Close()
	if errOp != nil {
		return errOp
	}

	fileTo, _ := os.Create(toPath)

	if limit == 0 || int(limit)+int(offset) > len(buf) {
		_, err = fileTo.Write(buf[offset:])
		if err != nil {
			log.Fatal("не записал")
		}
	} else {
		_, err = fileTo.Write(buf[offset : offset+limit])
		if err != nil {
			log.Fatal("не записал")
		}
	}

	errCl := fileTo.Close()
	if errCl != nil {
		return errCl
	}

	count := 100
	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.Finish()

	return nil
}
