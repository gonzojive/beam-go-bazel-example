package tfrecord

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

// RecordReaderOptions specify reader options for the tf record reader.
type RecordReaderOptions struct {
	CompressionType CompressionType
	// TODO: bufferSize?
	// TODO: zlibOptions?
}

// RecordReader implements a reader which can work on a queue of tf record files
// to extract the nested "Example" protobufs written using a matching tf record
// writer.
type RecordReader struct {
	queue   []string // file queue of files to read in
	options *RecordReaderOptions

	reader          *bufio.Reader
	recordsProduced int
}

// NewReader returns a new instance of a record reader which accepts a queue of
// files to read from.
func NewReader(queue []string, options *RecordReaderOptions) (*RecordReader, error) {
	return &RecordReader{
		queue:   queue,
		options: options,

		reader:          nil,
		recordsProduced: 0,
	}, nil
}

// NumRecordsProduced returns the number of records that this record reader has produced.
func (rr *RecordReader) NumRecordsProduced() int {
	return rr.recordsProduced
}

// readNextRecord will return the bytes that form the next successfully validated
// record found in the bytestream of the underlying reader.  If the reader returns
// an error, it is bubbled up (io.EOF is also an error, but just indicates that
// the reader is done, this is the caller's responsibility to handle).
// This method will return errors
func (rr *RecordReader) readNextRecord() ([]byte, error) {
	if rr.reader == nil {
		return nil, io.EOF
	}

	hbs, err := rr.reader.Peek(12)
	if err != nil {
		return nil, err
	}

	// Read the first 12 bytes into the entry structure and validate the length
	// field's CRC.
	var rec recordEntry
	if err := binary.Read(rr.reader, binary.LittleEndian, &rec.length); err != nil {
		return nil, err
	}
	if err := binary.Read(rr.reader, binary.LittleEndian, &rec.lengthCrc); err != nil {
		return nil, err
	}

	if MaskedCRC(hbs, 8) != rec.lengthCrc {
		return nil, errors.New("crc mismatch on record length")
	}

	// Read the length number of fields into the length array.
	offset := uint64(0)
	for offset < rec.length {
		left := uint64(4096)
		if rec.length-offset < left {
			left = rec.length - offset
		}
		chunk := make([]byte, left)
		n, err := rr.reader.Read(chunk)
		if err != nil {
			return nil, err
		}

		rec.data = append(rec.data, chunk[:n]...)
		offset += uint64(n)
	}

	if err := binary.Read(rr.reader, binary.LittleEndian, &rec.dataCrc); err != nil {
		return nil, err
	}
	if MaskedCRC(rec.data, int64(rec.length)) != rec.dataCrc {
		return nil, errors.New("crc mismatch on data")
	}

	return rec.data, nil
}

// ReadRecord checks the record reader for additional records and returns the next
// available one.  If the current reader returns an EOF, the queue is dequeued for
// another file to parse.  If the queue is empty, ReadRecord returns io.EOF to the
// caller.
func (rr *RecordReader) ReadRecord() ([]byte, error) {
	for {
		// If the reader is empty, and the queue has items - open the next item
		if rr.reader == nil && len(rr.queue) > 0 {
			var nextfp string
			nextfp, rr.queue = rr.queue[0], rr.queue[1:]
			f, err := os.Open(nextfp)
			if err != nil {
				return nil, err
			}
			rr.reader = bufio.NewReader(f)
		}

		// If the reader is nil - we are done, return io.EOF to signal we are done with
		// the last item in the queue.
		if rr.reader == nil {
			return nil, io.EOF
		}

		// Attempt to read an item off the reader, if we get an EOF - we need to try
		// and dequeue another work-item off the queue.
		bs, err := rr.readNextRecord()
		if err == io.EOF {
			rr.reader = nil
			continue
		} else if err == nil {
			rr.recordsProduced += 1
		}

		// Got a record (or an error), bubble up and done.
		return bs, err
	}
}
