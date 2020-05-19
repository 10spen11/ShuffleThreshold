package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	// the filepath
	filepath := "junkData.txt"

	// open a file for recording the test data
	file, err := os.Create(filepath)
	check(err)         // make sure nothing's wrong
	defer file.Close() // close the file eventually

	n := 5 // the number of contributors in the scheme
	t := 2 // the threshold

	// record the parameters of the test
	_, err = file.WriteString("Environment Creation:\n")
	check(err)
	_, err = file.WriteString("Using " + strconv.Itoa(n) + " voters with threshold " + strconv.Itoa(t) + "\n\n")
	check(err)

	start := time.Now() // start timer

	// create the environment for the tests
	shares := createThresholdShares(n, t) // this is where the bulk of the time is spent: overhead for creating the system
	// each user's share will contain the public key
	// Since these are all the same, we choose to use the copy at index 0 arbitrarlily
	publicKey := shares[0].Public()

	elapsed := time.Since(start)                        // end timer
	log.Printf("Environment Creation took %s", elapsed) // log the time
	_, err = file.WriteString(elapsed.String() + "\n")  // record the time in file
	check(err)

	// check each with an exponentially increasing number of ballots
	for ballotCount := 2; ballotCount < 100; ballotCount *= 2 {

		// state the number of ballots encrypted
		_, err = file.WriteString("\nEncrypting " + strconv.Itoa(ballotCount) + " ballots\n")
		check(err)

		for i := 0; i < 5; i++ { // do 5 tests

			start := time.Now() // start timer
			// doThresholdTest(ballotCount, n, t) // do test

			// generate messages
			messages, elGamal1, elGamal2 := generateMessageEncryptions(ballotCount, publicKey)

			// shuffle the messages
			elGamal1, elGamal2 = shuffleAndCheck(publicKey, elGamal1, elGamal2)

			// decrypt the messages, using the distributed shares
			decryptedMessages := decryptMessages(elGamal1, elGamal2, shares, t, n)

			// assures all decryptions are correct
			checkDecryption(messages, decryptedMessages)

			elapsed := time.Since(start)                       // end timer
			log.Printf("Decryption took %s", elapsed)          // log the time
			_, err = file.WriteString(elapsed.String() + "\n") // record the time in file
			check(err)
		}
	}

}
