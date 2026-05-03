package rnatranscription

func ToRNA(dna string) string {
	runes := []rune(dna)
    rnaMap := map[rune]rune{
        'G': 'C',
        'C': 'G',
        'T': 'A',
        'A': 'U',
    }
    for i, n := range runes{
        runes[i] = rnaMap[n]
    }
    return string(runes)
}
