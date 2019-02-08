// Package hhash implements a customizable, human-readable hashing engine.
package hhash

import (
	"encoding/binary"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/google/uuid"
)

// HHash uses xxhash to generate human-readable hash,
// with custom pattern and options for various parts of speech.
type HHash struct {
	// pattern is the hash pattern, set in InitPattern
	pattern string
	// regexp is a pointer to the compiled Regexp corresponding to the hash pattern, set in InitPattern
	regexp *regexp.Regexp
	// Adjectives is a slice of adjectives
	Adjectives []string
	// Adverbs is a slice of adverbs
	Adverbs []string
	// Adverbs is a slice of nouns
	Nouns []string
	// Verbs is a slice of verbs
	Verbs []string
	// VerbsPast is a slice of verbs in the past tense
	VerbsPast []string
	// VerbsGerund is a slice of verbs in the present participle
	VerbsGerund []string
	// AllowRepeats, if true, allows consecutive tokens to generate the same word (default: false)
	AllowRepeats bool
	// CalculateCollisionRate, if true, calculates the collision rate given the datasets, upon the first hash operation
	CalculateCollisionRate bool
	// numChoicesPerWord tracks the size of datasets considered in order to generate hash for some pattern; used to calculate collision rate
	numChoicesPerWord []uint
}

// WordType is a semantic type of words that a given token would generate.
// For example, the token "%A" represents a title-cased Adjective,
// and the "%v{G}" represents a a lowercased VerbGerund
type WordType int

const (
	// Adjective represents  adjectives.
	Adjective WordType = iota
	// Adverb represents adverbs.
	Adverb
	// Noun represents things and ideas but not persons or places.
	Noun
	// Verb represents actions in their present simple tense ("throw", "eat", "sing").
	Verb
	// VerbPast represents verbs in the past tense.
	VerbPast
	// VerbGerund represents verbs in the present participle tense.
	VerbGerund
)

// WordCase represents an enumeration of the categories in which a word may be cased.
type WordCase int

const (
	// Title casing uppercases the first letter.
	Title WordCase = iota
	// Lower casing lowercases the first letter.
	Lower
)

const patternRegex = `%(([a-zA-Z%])({([a-zA-Z\d]+)})?)`
const defaultPattern = `%A%V{G}%N`

// New instantiates an HHash with default datasets.
func New() *HHash {
	hasher := &HHash{}
	hasher.Adjectives = DefaultAdjectives
	hasher.Adverbs = DefaultAdverbs
	hasher.Nouns = DefaultNouns
	hasher.Verbs = DefaultVerbs
	hasher.VerbsGerund = DefaultVerbsGerund
	hasher.VerbsPast = DefaultVerbsPast
	return hasher
}

// NewWDefaultWithPattern instantiates an HHash with default datasets and the given pattern.
func NewWDefaultWithPattern(pattern string) *HHash {
	hasher := New()
	_ = hasher.InitPattern(pattern)
	return hasher
}

// NewDefault instantiates an HHash with default datasets and the default pattern.
func NewDefault() *HHash {
	return NewWDefaultWithPattern(defaultPattern)
}

// InitPattern initializes hhash with a new pattern of word tokens.
// TODO: examples
func (hasher *HHash) InitPattern(pattern string) error {
	regex, err := regexp.Compile(patternRegex)
	if err != nil {
		log.Printf("invalid regex %s: %s", patternRegex, err.Error())
		return err
	}

	log.Println("using pattern:", pattern)
	hasher.pattern = pattern
	hasher.regexp = regex
	hasher.numChoicesPerWord = make([]uint, 0)
	//randomHash := hasher.Random()
	//log.Println("random hash to calculate uniqueness probability:", randomHash)
	return nil
}

// hashForWordType produces a hash for the given word type and some previous hash.
// The previous hash is needed to avoid consecutively repetitive words.
// TODO: avoid repetitive words instead of consecutive ones only?
func (hasher *HHash) hashForWordType(wordType WordType, previousHash uint64) uint64 {
	var newHash uint64
	var hashBuffer []byte

	// re-hash the previous hash
	hashBuffer = make([]byte, 8) // 4*8=64 for uint64
	binary.LittleEndian.PutUint64(hashBuffer, previousHash)
	newHash = xxhash.Sum64(hashBuffer)

	if !hasher.AllowRepeats {
		wordsLength := uint64(len(hasher.WordsForType(wordType)))
		currentIndex := newHash % wordsLength
		previousIndex := previousHash % wordsLength
		timesRepeated := 0
		for currentIndex == previousIndex {
			timesRepeated++
			if timesRepeated > 1 {
				// for all three generatations of the seed (counting the initial) to be the same is extremely rare
				log.Printf(
					"previous hash (%d) and new hash (%d) map to the same index (%d) for the %d-th time. re-hashing...",
					previousHash, newHash, currentIndex, timesRepeated)
			}
			binary.LittleEndian.PutUint64(hashBuffer, newHash)
			newHash = xxhash.Sum64(hashBuffer)
			previousIndex = currentIndex
			currentIndex = newHash % wordsLength
		}
	}

	return newHash
}

// wordForToken returns a word mapped by the given hash in the given wordCase.
func (hasher *HHash) wordForToken(wordType WordType, wordCase WordCase, seed uint64) (string, uint64) {
	words := hasher.WordsForType(wordType)
	hash := hasher.hashForWordType(wordType, seed)
	index := hash % uint64(len(words))
	word := words[index]
	dataSetSize := uint(len(words))

	// numChoicesPerWord is initialized to empty slice
	if hasher.CalculateCollisionRate && hasher.numChoicesPerWord != nil {
		// track all the dataset sizes considered for calculating collision probability, but do it only once
		hasher.numChoicesPerWord = append(hasher.numChoicesPerWord, dataSetSize)
	}

	switch wordCase {
	case Lower:
		return strings.ToLower(word), hash
	default:
		return strings.Title(word), hash
	}
}

// WordsForType returns the dataset for the given word type.
func (hasher *HHash) WordsForType(wordType WordType) []string {
	switch wordType {
	case Adjective:
		return hasher.Adjectives
	case Adverb:
		return hasher.Adverbs
	case Noun:
		return hasher.Nouns
	case Verb:
		return hasher.Verbs
	case VerbPast:
		return hasher.VerbsPast
	case VerbGerund:
		return hasher.VerbsGerund
	default:
		log.Fatal("Unhandled WordType", wordType)
		return nil
	}
}

func (hasher *HHash) tokenReplacer(seed uint64) func(token string) string {
	return func(token string) string {
		// token-generated word
		var word string
		var wordType WordType
		var wordCase WordCase

		// parts comprise capture groups in patternRegex
		var parts []string
		// wordCode is the single character in token following the %, sans params
		var wordCode string
		var parameter string
		var err error

		// need to match for every token to find groups
		// see: https://github.com/golang/go/issues/5690
		// TODO: remove once issue is resolved
		parts = hasher.regexp.FindStringSubmatch(token)
		wordCode = parts[2]
		parameter = parts[4]
		wordType, wordCase, err = toWordType(wordCode, parameter)
		if err != nil {
			log.Printf("unable to determine word type for word code [%s] and parameter [%s]: %s\n", wordCode, parameter, err.Error())
			return token
		}

		// word is the hash value for the current token, and seed is updated for the next iteration of this inner method
		word, seed = hasher.wordForToken(wordType, wordCase, seed)
		if err != nil {
			log.Printf("unable to generate word for token [%s]: %s", token, err.Error())
			return token
		}

		return word
	}
}

// toWordType determines the word type and the word case given some word code and
func toWordType(wordCode, parameter string) (WordType, WordCase, error) {
	var wordCase WordCase
	lowerWordCode := strings.ToLower(wordCode)
	lowerParameter := strings.ToLower(parameter)

	if wordCode == lowerWordCode {
		wordCase = Lower
	} else {
		wordCase = Title
	}

	switch lowerWordCode {
	case "a":
		return Adverb, wordCase, nil
	case "j":
		return Adjective, wordCase, nil
	case "n":
		return Noun, wordCase, nil
	case "v":
		switch lowerParameter {
		case "p":
			return VerbPast, wordCase, nil
		case "g":
			return VerbGerund, wordCase, nil
		default:
			return Verb, wordCase, nil
		}
	}
	return -1, -1, fmt.Errorf("unable to determine WordType for flag [%s] and parameter [%s]", wordCode, parameter)
}

// HashUint returns a human-readable hash for the given seed.
func (hasher *HHash) HashUint(seed uint64) string {
	hashed := hasher.regexp.ReplaceAllStringFunc(hasher.pattern, hasher.tokenReplacer(seed))

	if hasher.CalculateCollisionRate && hasher.numChoicesPerWord != nil {
		var totalNumChoices uint = 1
		for _, numChoices := range hasher.numChoicesPerWord {
			totalNumChoices = totalNumChoices * numChoices
		}
		// TODO: do the calculation based on AllowRepeats
		log.Printf(
			"there is 1 in %d chance (%.16f%%) of hash collision given the current pattern (if allowing repeats)",
			totalNumChoices, float64(100)/float64(totalNumChoices))

		// null out hasher.numChoicesPerWord to stop calculating collision rate
		hasher.numChoicesPerWord = nil
	}
	return hashed
}

// HashString returns a human-readable hash for the given string.
func (hasher *HHash) HashString(s string) string {
	var seed = xxhash.Sum64String(s)
	return hasher.HashUint(seed)
}

// HashBytes returns a human-readable hash for the given string.
func (hasher *HHash) HashBytes(bytes []byte) string {
	var seed = xxhash.Sum64(bytes)
	return hasher.HashUint(seed)
}

// Random returns a human-readable hash based on a random UUID.
func (hasher *HHash) Random() string {
	randomUuid := uuid.New()
	return hasher.HashBytes(randomUuid[:])
}
