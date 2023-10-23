package palette_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"testing/iotest"

	"github.com/KlemensWinter/go-genie/palette"

	"github.com/stretchr/testify/assert"
)

func TestOpenPalette(t *testing.T) {
	pal, err := palette.Open("./testdata/50500_palette.txt")
	if assert.NoError(t, err) {
		assert.Len(t, pal, 256)
	}
}

func TestErrors(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		_, err := palette.Parse(iotest.ErrReader(io.EOF))
		assert.ErrorIs(t, err, io.EOF)
	})

	t.Run("invalid header", func(t *testing.T) {
		_, err := palette.Open("./testdata/invalid_header.txt")
		assert.ErrorIs(t, err, palette.ErrInvalidHeader)
	})
	t.Run("invalid version", func(t *testing.T) {
		_, err := palette.Open("./testdata/invalid_version.txt")
		assert.ErrorIs(t, err, palette.ErrInvalidVersion)
	})

	data, err := os.ReadFile("./testdata/50500_palette.txt")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("EOF in header", func(t *testing.T) {
		_, err := palette.Parse(
			io.LimitReader(bytes.NewReader(data), 2))
		assert.ErrorIs(t, err, io.EOF)
	})

	t.Run("EOF in version", func(t *testing.T) {
		_, err := palette.Parse(
			io.LimitReader(bytes.NewReader(data), 14))
		assert.ErrorIs(t, err, io.EOF)
	})

	t.Run("EOF in count", func(t *testing.T) {
		_, err := palette.Parse(
			io.LimitReader(bytes.NewReader(data), int64(len(palette.Header)+2+6)))
		assert.ErrorIs(t, err, io.EOF)
	})

	t.Run("invalid line", func(t *testing.T) {
		_, err := palette.Open("./testdata/invalid_line.txt")
		assert.Error(t, err)
	})

	t.Run("missing color", func(t *testing.T) {
		_, err := palette.Open("./testdata/missing_color.txt")
		assert.Error(t, err)
	})

	t.Run("chunked", func(t *testing.T) {
		_, err := palette.Parse(iotest.HalfReader(bytes.NewReader(data)))
		assert.NoError(t, err)
	})

	t.Run("Open not existing file", func(t *testing.T) {
		_, err := palette.Open("this_file_should_really_not_exist")
		assert.Error(t, err)
	})
}

func TestEncode(t *testing.T) {
	input, err := os.ReadFile("./testdata/50500_palette.txt")
	if err != nil {
		t.Fatal(err)
	}
	pal, err := palette.Parse(bytes.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	res, err := palette.Marshal(pal)
	if assert.NoError(t, err) {
		assert.Equal(t, input, res)
	}
}

func BenchmarkParse(b *testing.B) {
	data, err := os.ReadFile("./testdata/50500_palette.txt")
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		palette.Parse(bytes.NewReader(data))
	}
}
