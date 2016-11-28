package seq

import (
	"bytes"
	"errors"

	"github.com/kkaneda/genseq/fasta"
)

// overlap returns the length of the overlapping character if the
// given two sequences overlap. Returns -1 otherwise.
//
// The overlap condition is defined as follows. Let's c be the longest
// common substring satisfying the following conditions:
//
// - x = a::c (= c is x's suffix)
// - y = c::b (= c is y's prefix)
//
// x and y overlap if len(c) >= max(len(x)/2 + 1, len(y)/2 + 1).
//
// For example, if x is "CCTGCCGGAA", and y is "GCCGGAATAC",
// c is "GCCGGAA", x and y overlap since len(c) = 7 (>= 10/2 + 1 = 6).
func overlap(x, y string) int {
	// Iterate over [lower, upper) in the descending order where
	// lower = max(len(x)/2+1, len(y)/2+1) and
	// upper = min(len(x), len(y))
	lower := len(x)/2 + 1
	if l := len(y)/2 + 1; l > lower {
		lower = l
	}
	upper := len(x)
	if len(y) < upper {
		upper = len(y)
	}
	for n := upper - 1; n >= lower; n-- {
		matched := true
		// Check if x's suffix and y's prefix are equal
		// (length of both is n).
		for i := 0; i < n; i++ {
			if x[i+len(x)-n] != y[i] {
				matched = false
				break
			}
		}
		if matched {
			return n
		}
	}
	return -1
}

// edge is an edge of seqGraph.
type edge struct {
	// from is an index to the "from" node in the nodes array of the seqGraph.
	from int
	// to is an index to the "to" node in the nodes array of the seqGraph.
	to int
	// overlapLength is the length of the overlapping characters
	// between the from node and the to node.
	overlapLength int
}

// seqGraph defines a directed graph for sequences. Each vertex is a
// sequence and there is an edge from sequence X to another sequence Y
// if X and Y overlap.
//
// The graph is actually a list. For a given node, there should be at
// most one incoming edge and at most one outgoing edge.
type seqGraph struct {
	// nodes is a list of sequences.
	nodes []string
	// edges is a map of edges where keys are indices of "from" nodes.
	edges map[int]edge
	// revEdges is a map for reversed edges. Keys are indices of
	// "to" nodes. Values are indices of "from" nodes.
	revEdges map[int]int
}

// init initializes seqGraph by checking overlap for all possible combinations of
// sequences.
func (g *seqGraph) init(seqSet *fasta.SequenceSet) error {
	g.nodes = seqSet.Seqs
	g.edges = make(map[int]edge)
	g.revEdges = make(map[int]int)
	for i, seqI := range g.nodes {
		found := false
		for j, seqJ := range g.nodes {
			if i == j {
				continue
			}
			if length := overlap(seqI, seqJ); length > 0 {
				if found {
					return errors.New("there should be at most one overlapping sequence")
				}
				g.edges[i] = edge{i, j, length}
				g.revEdges[j] = i
				found = true
				// We could break the loop here but continue to make sure
				// seqI has the exactly one overlapping sequence.
			}
		}
	}
	return nil
}

// head returns the head of the graph. Returns -1 if no head node is found.
func (g *seqGraph) head() int {
	for i := range g.nodes {
		if _, ok := g.revEdges[i]; !ok {
			return i
		}
	}
	return -1
}

// traverse traverses the graph from the head node to the tail node
// and returns the output sequence.
func (g *seqGraph) traverse() (string, error) {
	curr := g.head()
	if curr < 0 {
		return "", errors.New("no head node found")
	}
	buf := new(bytes.Buffer)
	buf.WriteString(g.nodes[curr])
	seqCount := 0
	for {
		seqCount++
		edge, ok := g.edges[curr]
		if !ok {
			break
		}
		curr = edge.to
		buf.WriteString(g.nodes[curr][edge.overlapLength:])
	}

	// Check if we visited the all the sequences.
	if seqCount != len(g.nodes) {
		return "", errors.New("all sequences must be traversed")
	}
	return buf.String(), nil
}

// Run runs the sequence algorithm and returns its output.
func Run(seqSet *fasta.SequenceSet) (string, error) {
	var g seqGraph
	g.init(seqSet)
	return g.traverse()
}
