# run-match
GO client and main for IDSL code assessment

The run-match Go application is the main for the biometric match project.

Here the MdTF repository is cloned into memory and the ~/face images (*.png) are copied to /tmp/images directory.  (Upon reboot, /tmp is cleared.)  Arraya of the file names are utilized to run the MatchFiles() via the api client to match 2 files and retrieves a score.  Each image file is matched to the other image file.

Once all images are processed, the api client's GetAllMatchScores() is used to obtain a list of all image comparisions and their scores.  A simple table is output to the terminal to show FILENAM1, FILENAME2, and MATCHSCORE.

To Run: $go run .

TODO test files for functions used in runmatch.go.  I have run out of time.
