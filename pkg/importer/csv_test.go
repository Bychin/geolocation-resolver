package importer

import (
	"testing"
)

func TestCSVImporter(t *testing.T) {
	const pathPrefix = "testdata/"

	for i, testcase := range []struct {
		filename      string
		validLinesNum int
		mustFail      bool
	}{
		{
			filename:      "data_dump_valid.csv",
			validLinesNum: 6,
			mustFail:      false,
		},
		{
			filename:      "data_dump_with_duplicate_lines.csv",
			validLinesNum: 1, // must ignore duplicate by IP
			mustFail:      false,
		},
		{
			filename:      "data_dump_with_invalid_lines.csv",
			validLinesNum: 3, // must ignore invalid lines
			mustFail:      false,
		},
		{
			filename:      "data_dump_with_line_with_extra_column.csv",
			validLinesNum: 3, // must ignore extra column
			mustFail:      false,
		},
		{
			filename:      "data_dump_without_mystery_value_column.csv",
			validLinesNum: 2, // must ignore extra column
			mustFail:      false,
		},
		{
			filename:      "data_dump_with_invalid_line_content.csv",
			validLinesNum: 0, // must ignore line with invalid content
			mustFail:      false,
		},
		{
			filename:      "data_dump_with_unescaped_sequence.csv",
			validLinesNum: 1,
			mustFail:      true, // invalid sequence on last line
		},
		{
			filename:      "wrong_path",
			validLinesNum: 0,
			mustFail:      true, // invalid sequence on last line
		},
	} {
		importer := NewCSVImporter(pathPrefix+testcase.filename, ',', true)

		done := make(chan struct{})
		go func() {
			for range importer.Entries() {
				// pass
			}
			close(done)
		}()

		stats, importErr := importer.Exec()
		<-done

		if stats.ParsedLines != testcase.validLinesNum {
			t.Errorf("#%d: got unexpected amount of valid lines in CSV, have: %d, want: %d",
				i, stats.ParsedLines, testcase.validLinesNum)
		}
		if importErr != nil && !testcase.mustFail {
			t.Errorf("#%d: got error when parsing CSV: %s", i, importErr)
			continue
		}
		if importErr == nil && testcase.mustFail {
			t.Errorf("#%d: expected to fail, but didn't", i)
		}
	}
}
