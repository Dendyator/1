package main

import (
	"errors"
	"github.com/cheggaaa/pb/v3"
	"io"
	"log"
	"os"
	"time"
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

	if limit == 0 || int(limit) > len(buf) {
		fileTo, _ := os.Create(toPath)
		_, err = fileTo.Write(buf[offset:])
		if err != nil {
			log.Fatal("не записал")
		}
		e := fileTo.Close()
		if e != nil {
			return e
		}
	} else {
		fileTo, _ := os.Create(toPath)
		_, err = fileTo.Write(buf[offset : offset+limit])
		if err != nil {
			log.Fatal("не записал")
		}
		e := fileTo.Close()
		if e != nil {
			return e
		}
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
