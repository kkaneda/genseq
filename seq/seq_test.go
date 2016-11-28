package seq

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/kkaneda/genseq/fasta"
)

// TestOverlap tests overlap.
func TestOverlap(t *testing.T) {
	type testCase struct {
		seq1, seq2 string
		expected   int
	}
	testCases := []testCase{
		{"ATC", "TCG", 2},
		{"CAAAA", "AAAAC", 4},
		{"ATTAGACCTG", "AGACCTGCCG", 7},
		{"AAC", "AAC", -1},
		{"A", "ATC", -1},
		{"ATTAGACCTG", "CCTGCCGGAA", -1},
	}
	for i, testCase := range testCases {
		if ret := overlap(testCase.seq1, testCase.seq2); ret != testCase.expected {
			t.Errorf("%d: expected %d, but got %d", i, testCase.expected, ret)
		}
	}
}

// BenchmarkOverlap runs the benchmark for overlap.
func BenchmarkOverlap(b *testing.B) {
	seqSet, err := fasta.LoadFromFile("../testdata/contest.txt")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		for _, seqI := range seqSet.Seqs {
			for _, seqJ := range seqSet.Seqs {
				overlap(seqI, seqJ)
			}
		}
	}
}

// BenchmarkOverlap2 runs the benchmark for overlap2.
func BenchmarkOverlap2(b *testing.B) {
	seqSet, err := fasta.LoadFromFile("../testdata/contest.txt")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		for _, seqI := range seqSet.Seqs {
			for _, seqJ := range seqSet.Seqs {
				overlap2(seqI, seqJ)
			}
		}
	}
}

// TestRun_Small tests the sequence algorithm for a small input file.
func TestRun_Small(t *testing.T) {
	seqSet, err := fasta.LoadFromFile("../testdata/small.txt")
	if err != nil {
		t.Fatal(err)
	}
	output, err := Run(seqSet)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "ATTAGACCTGCCGGAATAC"; strings.Compare(output, expected) != 0 {
		t.Errorf("expected %s, but got %s", expected, output)
	}
}

// TestRun_Contest tests the sequence algorithm for the contest input.
func TestRun_Contest(t *testing.T) {
	seqSet, err := fasta.LoadFromFile("../testdata/contest.txt")
	if err != nil {
		t.Fatal(err)
	}
	output, err := Run(seqSet)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := ioutil.ReadFile("../testdata/contest_output.txt")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Compare(output, string(expected)) != 0 {
		t.Errorf("expected %s, but got %s", len(expected), len(output))
	}
}

// BenchmarkRun runs the benchmark of the sequence algorithm.
func BenchmarkRun(b *testing.B) {
	seqSet, err := fasta.LoadFromFile("../testdata/contest.txt")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		if _, err := Run(seqSet); err != nil {
			b.Fatal(err)
		}
	}
}
