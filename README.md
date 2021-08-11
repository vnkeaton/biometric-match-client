# biometric-match-client

The biometric-match-client is a Go api client for the java REST service biometric-match.

The client provides the ability to match images via the MatchFiles() function.  The function accepts 2 image files for comparison where it calls the ~/match endpoint of the REST api and returns the match score information that consists of the files' information and the comparison match score.

The client also provides a function called GetAllMatchScores() that will consume all match scores that have been posted.  The function calls the ~/biometric/matchscore/downloadFile/all and returns an array that consist of two files' informations and it's matching score.

The biometric-match-client Go api client is packaged as matchclient and is imported in the main run-match Go application.

There is an accompaning test file, matchclient_test.go.  The tests included are for the Hello() fucntion that hits the ~/biometric/hello endpoint, the Matchfiles() that hits the ~/biometric/image/match endpoint, and the GetAllMatchScores() that hits the ~biometric/matchscore/downloadFile/all endpoint.

To run the test: $go test
