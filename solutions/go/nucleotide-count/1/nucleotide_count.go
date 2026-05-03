package nucleotidecount
import "fmt"
// Histogram is a mapping from nucleotide to its count in given DNA.
// Choose a suitable data type.
// Start by uncommenting the following line:
// type Histogram ...
type Histogram map[rune]int
// DNA is a list of nucleotides. Choose a suitable data type.
// Start by uncommenting the following line:
// type DNA ...
type DNA string
// Counts generates a histogram of valid nucleotides in the given DNA.
// Returns an error if d contains an invalid nucleotide.
//
// Counts is a method on the DNA type. A method is a function with a special receiver argument.
// The receiver appears in its own argument list between the func keyword and the method name.
// Here, the Counts method has a receiver of type DNA named d.
func (d DNA) Counts() (Histogram, error) {
	counts := Histogram {
        'A': 0,
        'C': 0,
        'G': 0,
        'T': 0,
    }
    for _, nucl := range d{
        switch nucl {
            case 'A', 'C', 'G', 'T':
            counts[nucl]++
        	default:
            return nil, fmt.Errorf("знайдено невалідний нуклеотид: %c", nucl)
        }
    }
	return counts, nil
}
