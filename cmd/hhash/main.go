package main

import (
	"bufio"
	"fmt"
	"hhash"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/integrii/flaggy"
)

const defaultPattern = "%j_%n"

func main() {
	start := time.Now()
	pattern := defaultPattern
	verbose := false
	allowRepeats := false
	toHash := ""

	flaggy.String(&pattern, "p", "pattern", fmt.Sprintf("hash pattern consisting of static characters and word tokens (default: \"%s\")", defaultPattern))
	flaggy.String(&toHash, "s", "toHash", "string to hash according to the pattern (default: random)")
	flaggy.Bool(&verbose, "v", "verbose", "enable verbose logging (default: false)")
	flaggy.Bool(&allowRepeats, "r", "repetition", "allow same word to be generated from consecutive tokens that are the same (e.g., \"%N%N\")")
	flaggy.Parse()

	if !verbose {
		// causes only the final hash to be printed to stdout
		log.SetOutput(ioutil.Discard)
	}

	// init hhash
	hasher := hhash.New()
	hasher.AllowRepeats = allowRepeats
	hasher.CalculateCollisionRate = verbose
	err := hasher.InitPattern(pattern)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("invalid pattern %s: %s\n", pattern, err.Error()))
		os.Exit(1)
	}

	// check for pipe
	var hashed string
	info, _ := os.Stdin.Stat()
	if info.Size() > 0 {
		if toHash != "" {
			_, _ = os.Stderr.WriteString("-s must not be set when piping to hhash\n")
			os.Exit(1)
		}
		hashed = hashPipedStream(hasher)
	} else {
		// if no string is provided to hash, a random string is generated based on uuid
		if toHash == "" {
			randomUUID := uuid.New()
			toHash = string(randomUUID[:])
		}
		hashed = hasher.HashString(toHash)
	}

	elapsed := time.Since(start)
	log.Printf("elapsed time %s", elapsed)
	fmt.Println(hashed)
}

func hashPipedStream(hasher *hhash.HHash) string {
	reader := bufio.NewReader(os.Stdin)
	lastHashValue := ""
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		lastHashValue = hasher.HashString(lastHashValue + line)
	}
	return lastHashValue
}
