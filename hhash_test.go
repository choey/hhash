package hhash

import (
	"log"
	"testing"
	"time"

	"github.com/agnivade/levenshtein"
	"github.com/google/uuid"
)

const allInclusivePattern = "%v{G}%N-%a%V-%J%V{P}"

func TestKnownHash(t *testing.T) {
	hhash := NewDefault()
	startTime := time.Now()
	hash1 := hhash.HashString("pumba & mumble")
	hash2 := hhash.HashString("pumba + mumble")
	distance := levenshtein.ComputeDistance(hash1, hash2)
	minDistance := (len(hash1) + len(hash2)) / 3
	t.Logf("hash1: %s, hash2: %s\n", hash1, hash2)
	t.Log("time to hash:", time.Since(startTime))
	if distance < minDistance {
		t.Logf("levenshtein distance between [%s] and [%s] was %d but should be greater than %d", hash1, hash2, distance, minDistance)
		t.Fail()
	}
}

func TestRandomHashes(t *testing.T) {
	startTime := time.Now()
	hhash := NewDefault()
	err := hhash.InitPattern(allInclusivePattern)
	log.Println("time took to init hhash:", time.Since(startTime))
	if err != nil {
		log.Fatal(err)
	}

	// pre-calculate seeds
	const iterations = 100000
	var seeds [iterations][]byte
	for i := 0; i < iterations; i = i + 1 {
		uuid := uuid.New()
		seeds[i] = uuid[:]
	}

	// hash all the pre-calculated seeds
	startTime = time.Now()
	for i := 0; i < iterations; i = i + 1 {
		hhash.HashBytes(seeds[i])
	}

	t.Logf("time to generate %d GUID-based hashes: %s", iterations, time.Since(startTime))
}
