package fasta

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/kkaneda/genseq/testutils"
)

// TestLoadFromFile_Small tests LoadFromFile with the small test file.
func TestLoadFromFile_Small(t *testing.T) {
	seqSet, err := LoadFromFile("../testdata/small.txt")
	if err != nil {
		t.Fatal(err)
	}
	expectedSeqs := []string{
		"ATTAGACCTG",
		"CCTGCCGGAA",
		"AGACCTGCCG",
		"GCCGGAATAC",
	}
	if len(seqSet.Seqs) != len(expectedSeqs) {
		t.Fatalf("expected %d seqs, but got %d", len(expectedSeqs), len(seqSet.Seqs))
	}
	for i, seq := range seqSet.Seqs {
		if expected := expectedSeqs[i]; strings.Compare(seq, expected) != 0 {
			t.Errorf("expected %s, but got %s", expected, seq)
		}
	}
}

// TestLoadFromFile_Contest tests LoadFromFile with the contest file.
func TestLoadFromFile_Contest(t *testing.T) {
	seqSet, err := LoadFromFile("../testdata/contest.txt")
	if err != nil {
		t.Fatal(err)
	}

	if expectedLen := 50; len(seqSet.Seqs) != expectedLen {
		t.Fatalf("expected %d seqs, but got %d", expectedLen, len(seqSet.Seqs))
	}
}

// TestLoadFromFile_Malformed tests LoadFromFile with the contest file.
func TestLoadFromFile_Malformed(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	file, err := ioutil.TempFile(dir, "test.txt")
	if err != nil {
		t.Fatal(err)
	}
	content := []byte("ATCG\n")
	if _, err := file.Write(content); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFromFile(file.Name()); !testutils.IsError(err, "sequence started without a descriptor line: line no: 1") {
		t.Errorf("unexpected error %v", err)
	}
}

// TestValidate_TooManySeqs verifies that the validation fails when
// there are too many sequences.
func TestValidate_TooManySeqs(t *testing.T) {
	var seqSet SequenceSet
	seqSet.Seqs = make([]string, 0, 0)
	if err := validate(&seqSet); !testutils.IsError(err, "no sequence is found") {
		t.Errorf("unexpected error %v", err)
	}

	seqSet.Seqs = make([]string, maxSequencePerFile+1)
	if err := validate(&seqSet); !testutils.IsError(err, "too many sequences are found: .*") {
		t.Errorf("unexpected error %v", err)
	}

	seqSet.Seqs = []string{""}
	if err := validate(&seqSet); !testutils.IsError(err, "sequence must have at least one character") {
		t.Errorf("unexpected error %v", err)
	}

	seq := make([]byte, maxLengthPerSeq+1)
	seqSet.Seqs = []string{string(seq)}
	if err := validate(&seqSet); !testutils.IsError(err, "too long sequence found: .*") {
		t.Errorf("unexpected error %v", err)
	}

	seqSet.Seqs = []string{"D"}
	if err := validate(&seqSet); !testutils.IsError(err, "unexpected character D") {
		t.Errorf("unexpected error %v", err)
	}
}
