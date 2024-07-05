package bag

import "math"

func New(c Config) *Bag {
	var b Bag
	b.c = c
	b.vocabByLabel = map[string]Vocabulary{}
	b.countByLabel = map[string]int{}
	// Fill unset values as default
	b.c.fill()
	return &b
}

type Bag struct {
	// Configuration values
	c Config
	// Vocabulary sets by label
	vocabByLabel map[string]Vocabulary
	// Count of trained documents by label
	countByLabel map[string]int
	// Total count of trained documents
	totalCount int
}

func (b *Bag) GetResults(in string) (r Results) {
	// Convert inbound data to NGrams
	ns := toNGrams(in, b.c.NGramSize)
	// Initialize results with the same size as the current number of vocabulary labels
	r = make(Results, len(b.vocabByLabel))
	// Iterate through vocabulary sets by label
	for label, vocab := range b.vocabByLabel {
		// Set probability value for iterating label
		r[label] = b.getProbability(ns, label, vocab)
	}

	return
}

func (b *Bag) Train(in, label string) {
	// Convert inbound data to a slice of NGrams
	ns := toNGrams(in, b.c.NGramSize)
	// Get vocabulary for a provided label, if the vocabulary doesn't exist, it will be created)
	v := b.getOrCreateVocabulary(label)
	// Iterate through NGrams
	for _, n := range ns {
		// Increment the vocabulary value for the current NGram
		v[n.String()]++
	}

	// Increment count of trained documents for the provided label
	b.countByLabel[label]++
	// Increment total count of trained documents
	b.totalCount++
}

func (b *Bag) getProbability(ns []NGram, label string, vocab Vocabulary) (probability float64) {
	// Set initial probability value as the prior probability value
	probability = b.getPriorProbability(label)
	// Get the current counts by label (to be used by Laplace smoothing during for-loop)
	countsByLabel := float64(b.countByLabel[label] + len(vocab))
	// Iterate through NGrams
	for _, n := range ns {
		// Utilize Laplace smoothing to improve our results when an ngram isn't found within the trained dataset
		// Likelihood with Laplace smoothing
		count := float64(vocab[n.String()] + b.c.SmoothingParameter)
		// Add logarithmic result of count (plus )
		probability += math.Log(count / countsByLabel)
	}

	return
}

func (b *Bag) getPriorProbability(label string) (probability float64) {
	count := float64(b.countByLabel[label])
	total := float64(b.totalCount)
	// Get the logarithmic value of count divided by total count
	return math.Log(count / total)
}

func (b *Bag) getOrCreateVocabulary(label string) (v Vocabulary) {
	var ok bool
	v, ok = b.vocabByLabel[label]
	// Check if vocabulary set does not exist for the provided label
	if !ok {
		// Create new vocabulary set
		v = make(Vocabulary)
		// Set vocabulary set by label as newly created value
		b.vocabByLabel[label] = v
	}

	return
}
