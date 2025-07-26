package isaac

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestSeededMemory(t *testing.T) {
	// note:
	//   this test compares the internal results slice with values from the official reference implementation after being
	//   initialized with the same seed values.

	// open the applicable testdata file for reading
	f, err := os.Open("testdata/keytest.txt")
	if err != nil {
		t.Fatalf("failed to open keytest.txt: %v", err)
	}
	defer f.Close()

	// read the entire test data file into memory
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read keytest.txt: %v", err)
	}
	expected := strings.Split(string(b), "\n")[:32]

	// transform the required seed string into a usable seed slice
	seed, err := TransformSeed("This is <i>not</i> the right mytext.")
	if err != nil {
		t.Fatalf("failed to transform seed into slice: %v", err)
	}

	// create the pseudorandom number generate with the required seed values
	r := NewRand(seed...)

	// iterate one time for each entry in the results slice
	var actual bytes.Buffer
	for i := uint32(0); i < Size; i++ {
		// add the value to the actual output
		_, _ = fmt.Fprintf(&actual, "%08x ", r.results[i])

		// check if a line break needs to be added
		if i&7 == 7 {
			_, _ = fmt.Fprintln(&actual)
		}
	}

	// loop through the expected output
	output := strings.Split(actual.String(), "\n")
	for i, value := range expected {
		// check if the value matches in the actual output
		if value != output[i] {
			t.Fatalf("value mismatch on line %d: expected %q, actual %q", i+1, value, output[i])
		}
	}
}

func TestTransformSeed(t *testing.T) {
	var seedString string
	var seedSlice []uint32

	// verify the sample seed used in the official tests is transformed correctly
	seedString = "This is <i>not</i> the right mytext."
	seedSlice = []uint32{
		0x73696854,
		0x20736920,
		0x6e3e693c,
		0x2f3c746f,
		0x74203e69,
		0x72206568,
		0x74686769,
		0x74796D20,
		0x2E747865,
	}

	// transform the seed string into a seed slice
	seed, err := TransformSeed(seedString)
	if err != nil {
		t.Fatalf("failed to transform seed string %q: %v", seedString, err)
	}

	// verify the seed lengths are identical
	if len(seedSlice) != len(seed) {
		t.Fatalf("incorrectly transformed seed string %q mismatched lengths: expected %d, got %d", seedString, len(seedSlice), len(seed))
	}

	// loop through the values in the seed slice
	for i := range seedSlice {
		// verify the value matches the expected value
		if seedSlice[i] != seed[i] {
			t.Fatalf("incorrectly transformed seed string %q at index %d: expected 0x%08x, got 0x%08x", seedString, i, seedSlice[i], seed[i])
		}
	}

	// remove two characters from the sample seed to ensure it requires padding
	seedString = "This is <i>not</i> the right mytex"
	seedSlice = []uint32{
		0x73696854,
		0x20736920,
		0x6e3e693c,
		0x2f3c746f,
		0x74203e69,
		0x72206568,
		0x74686769,
		0x74796D20,
		0x00007865,
	}

	// transform the seed string into a seed slice
	seed, err = TransformSeed(seedString)
	if err != nil {
		t.Fatalf("failed to transform seed string %q: %v", seedString, err)
	}

	// verify the seed lengths are identical
	if len(seedSlice) != len(seed) {
		t.Fatalf("incorrectly transformed seed string %q mismatched lengths: expected %d, got %d", seedString, len(seedSlice), len(seed))
	}

	// loop through the values in the seed slice
	for i := range seedSlice {
		// verify the value matches the expected value
		if seedSlice[i] != seed[i] {
			t.Fatalf("incorrectly transformed seed string %q at index %d: expected 0x%08x, got 0x%08x", seedString, i, seedSlice[i], seed[i])
		}
	}
}

func TestWithSeed(t *testing.T) {
	// note:
	//   this test compares the library output against the `randseed.txt` file that can be found on:
	//   https://www.burtleburtle.net/bob/rand/isaacafa.html

	// open the applicable testdata file for reading
	f, err := os.Open("testdata/randseed.txt")
	if err != nil {
		t.Fatalf("failed to open randseed.txt: %v", err)
	}
	defer f.Close()

	// read the entire test data file into memory
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read randseed.txt: %v", err)
	}
	expected := strings.Split(string(b), "\n")[:320]

	// transform the required seed string into a usable seed slice
	seed, err := TransformSeed("This is <i>not</i> the right mytext.")
	if err != nil {
		t.Fatalf("failed to transform seed into slice: %v", err)
	}

	// create the pseudorandom number generate with the required seed values
	r := NewRand(seed...)

	// iterate exactly 10 times
	var actual bytes.Buffer
	var k uint32
	for i := uint32(0); i < 10; i++ {
		// iterate one time for each entry in the results slice
		for j := uint32(0); j < Size; j++ {
			// add the value to the actual output
			_, _ = fmt.Fprintf(&actual, "%08x ", r.Uint32())

			// check if a line break needs to be added
			k += 1
			if k == 8 {
				k = 0
				_, _ = fmt.Fprintln(&actual)
			}
		}
	}

	// loop through the expected output
	output := strings.Split(actual.String(), "\n")
	for i, value := range expected {
		// check if the value matches in the actual output
		if value != output[i] {
			t.Fatalf("value mismatch on line %d: expected %q, actual %q", i+1, value, output[i])
		}
	}
}

func TestZeroSeed(t *testing.T) {
	// note:
	//   this test compares the library output against the `randvect.txt` file that can be found on:
	//   https://www.burtleburtle.net/bob/rand/isaacafa.html

	// open the applicable testdata file for reading
	f, err := os.Open("testdata/randvect.txt")
	if err != nil {
		t.Fatalf("failed to open randvect.txt: %v", err)
	}
	defer f.Close()

	// read the entire test data file into memory
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read randvect.txt: %v", err)
	}
	expected := strings.Split(string(b), "\n")[:64]

	// create the pseudorandom number generate with no seed (all zeros, by default)
	r := NewRand()

	// iterate exactly twice
	var actual bytes.Buffer
	for i := uint32(0); i < 2; i++ {
		// generate a new set of pseudorandom numbers
		r.isaac()

		// iterate one time for each entry in the results slice
		for j := uint32(0); j < Size; j++ {
			// add the value to the actual output
			_, _ = fmt.Fprintf(&actual, "%08x", r.results[j])

			// check if a line break needs to be added
			if j&7 == 7 {
				_, _ = fmt.Fprintln(&actual)
			}
		}
	}

	// loop through the expected output
	output := strings.Split(actual.String(), "\n")
	for i, value := range expected {
		// check if the value matches in the actual output
		if value != output[i] {
			t.Fatalf("value mismatch on line %d: expected %q, actual %q", i+1, value, output[i])
		}
	}
}
