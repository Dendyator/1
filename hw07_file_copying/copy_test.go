package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const toPath = "D:/Coding/Projects/NinjaProject/Uroki/Les1/output.txt"

func TestCopy(t *testing.T) {
	t.Run("bad file", func(t *testing.T) {
		fromPath := "D:/Coding/Projects/NinjaProject/Uroki/Les1/input1.txt"

		offset = 6000
		limit = 0

		err := Copy(fromPath, toPath, offset, limit)
		if err != nil {
			fmt.Println(err)
		}

		require.EqualError(t, err, "unsupported file")
	})
}

func TestCopy2(t *testing.T) {
	t.Run("big offset", func(t *testing.T) {
		fromPath := "D:/Coding/Projects/NinjaProject/Uroki/Les1/input.txt"
		offset = 8000
		limit = 0

		err := Copy(fromPath, toPath, offset, limit)
		if err != nil {
			fmt.Println(err)
		}

		require.EqualError(t, err, "offset exceeds file size")
	})
}

func TestCopy3(t *testing.T) {
	t.Run("len new file", func(t *testing.T) {
		fromPath := "D:/Coding/Projects/NinjaProject/Uroki/Les1/input.txt"
		offset = 3000
		limit = 1000

		err := Copy(fromPath, toPath, offset, limit)
		if err != nil {
			fmt.Println(err)
		}

		file, err := os.Open(toPath)
		if err != nil {
			log.Fatal("Нет файла")
		}

		buf, er := io.ReadAll(file)
		if er != nil {
			log.Fatal("не прочитал")
		}

		require.Len(t, buf, 1000)
		require.NoError(t, err)
	})
}
