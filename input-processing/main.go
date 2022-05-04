package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sync"
)

type LineReader struct {
	stream <-chan byte
}

func NewLineReader(stream <-chan byte) *LineReader {
	return &LineReader{
		stream: stream,
	}
}

func (s *LineReader) ReadRune() (r rune, size int, err error) {
	r = rune(<-s.stream)
	if r == '\n' {
		return 0, 0, io.EOF
	} else {
		return r, 4, nil
	}
}

type LineFileBuffer struct {
	file           *os.File
	bufferedWriter *bufio.Writer
	lineSize       int64
}

func NewLineFileBuffer() *LineFileBuffer {
	if _, err := os.Stat("line-buffer.txt"); err == nil {
		os.Remove("line-buffer.txt")

	} else if !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	file, err := os.Create("line-buffer.txt")

	if err != nil {
		panic(err)
	}

	bufferedWriter := bufio.NewWriter(file)

	return &LineFileBuffer{
		file:           file,
		bufferedWriter: bufferedWriter,
	}
}

func (b *LineFileBuffer) Read(p []byte) (n int, err error) {
	return b.file.Read(p)
}

func (b *LineFileBuffer) Size() int64 {
	return b.lineSize
}

func (b *LineFileBuffer) WriteByte(p byte) error {
	b.lineSize++
	return b.bufferedWriter.WriteByte(p)
}

func (b *LineFileBuffer) Flush() error {
	return b.bufferedWriter.Flush()
}

func (b *LineFileBuffer) ResetPosition() (err error) {
	_, err = b.file.Seek(0, 0)

	return err
}

func (b *LineFileBuffer) Reset() (err error) {
	b.lineSize = 0
	b.bufferedWriter.Reset(b.file)
	return b.ResetPosition()
}

func (b *LineFileBuffer) Close() error {
	err := b.file.Close()

	if err != nil {
		return err
	}

	return os.Remove("line-buffer.txt")
}

func scan(reader *bufio.Reader, regex string, lineScannedCallback func(reader io.Reader, lineSize int64, matched bool)) {
	parsedRegex := regexp.MustCompile(regex)

	regexStream := make(chan byte)

	lineStream := make(chan byte)

	lineReader := NewLineReader(regexStream)

	inputBuffer := make([]byte, reader.Size())

	lineBuffer := NewLineFileBuffer()

	matchedReader := make(chan bool)

	var taskGroup sync.WaitGroup

	go func() {
		for {
			matchedReader <- parsedRegex.MatchReader(lineReader)
		}
	}()

	taskGroup.Add(1)

	go func() {
		for {

			matchFound := false

		inputStreamReceiver:
			for b := range lineStream {
				err := lineBuffer.WriteByte(b)

				if err != nil {
					panic(err)
				}

				select {
				case regexStream <- b:
					if b == '\n' {
						matchFound = <-matchedReader
						break inputStreamReceiver
					}

				case matchFound = <-matchedReader:
					for {
						if b == '\n' {
							break inputStreamReceiver
						}

						b = <-lineStream

						err := lineBuffer.WriteByte(b)

						if err != nil {
							panic(err)
						}
					}
				}
			}

			if lineBuffer.Size() == 0 {
				break
			}

			lineBuffer.Flush()

			err := lineBuffer.ResetPosition()

			if err != nil {
				panic(err)
			}

			if matchFound {
				lineScannedCallback(lineBuffer, lineBuffer.Size(), true)
			} else {
				lineScannedCallback(lineBuffer, lineBuffer.Size(), false)
			}

			lineBuffer.Reset()
		}

		taskGroup.Done()
	}()

	for {
		inputChunkSize, err := reader.Read(inputBuffer)

		for i := 0; i < inputChunkSize; i++ {
			lineStream <- inputBuffer[i]
		}

		if err == io.EOF {
			lineStream <- '\n'
			close(lineStream)
			break
		}
	}

	taskGroup.Wait()

	lineBuffer.Close()
}

func main() {
	fmt.Println("SP// Backend Developer Test - Input Processing")
	fmt.Println()

	// Read STDIN into a new buffered reader
	reader := bufio.NewReader(os.Stdin)

	// TODO: Look for lines in the STDIN reader that contain "error" and output them.

	scan(reader, "^(?:error|.*[: .]error)(?:(?:[: .]+.*)$|$)", func(lineReader io.Reader, lineSize int64, matched bool) {
		if matched {
			io.CopyN(os.Stdout, lineReader, lineSize)
		}
	})
}
