package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	fromUID int
	toUID   int

	fromGID int
	toGID   int

	fromUser string
	toUser   string

	fromGroup string
	toGroup   string

	verbose bool
)

func init() {
	flag.IntVar(&fromUID, "from-uid", fromUID, "Rewrite UID from UID")
	flag.IntVar(&toUID, "to-uid", toUID, "Rewrite UID to  UID")

	flag.IntVar(&fromGID, "from-gid", fromGID, "Rewrite GID from GID")
	flag.IntVar(&toGID, "to-gid", toGID, "Rewrite this GID to GID")

	flag.StringVar(&fromUser, "from-user", fromUser, "Rewrite username FROM username")
	flag.StringVar(&toUser, "to-user", toUser, "Rewrite username FROM username")

	flag.StringVar(&fromGroup, "from-group", fromUser, "Rewrite groupname FROM groupname")
	flag.StringVar(&toGroup, "to-group", toUser, "Rewrite groupname FROM groupname")

	flag.BoolVar(&verbose, "verbose", verbose, "Be verbose about changes")
}

func main() {
	flag.Parse()

	var potentialChanges = 0
	if fromUID != toUID {
		if verbose {
			log.Printf("Rewriting UID: %d->%d", fromUID, toUID)
		}
		potentialChanges++
	}
	if fromGID != toGID {
		if verbose {
			log.Printf("Rewriting GID: %d->%d", fromGID, toGID)
		}
		potentialChanges++
	}
	if fromUser != toUser {
		if verbose {
			log.Printf("Rewriting user name: %s->%s", fromUser, toUser)
		}
		potentialChanges++
	}
	if fromGroup != toGroup {
		if verbose {
			log.Printf("Rewriting group name: %s->%s", fromGroup, toGroup)
		}
		potentialChanges++
	}

	if potentialChanges == 0 {
		if verbose {
			log.Println("Effectively no changes requested. Passing data through untouched")
		}
		_, err := io.Copy(os.Stdout, os.Stdin)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	tarInput := tar.NewReader(os.Stdin)
	tarOutput := tar.NewWriter(os.Stdout)

	effectiveChanges := 0
	entryChanges := 0
	totalEntries := 0

	for {
		header, err := tarInput.Next()
		if err != nil {
			if err != tar.ErrInsecurePath {
				if err == io.EOF {
					break
				}
			} else {
				panic(err)
			}
		}

		totalEntries++
		verboseOut := []string{}

		if fromUID == header.Uid {
			header.Uid = toUID
			if verbose {
				effectiveChanges++
				verboseOut = append(verboseOut, fmt.Sprintf("UID:%d:%d", fromUID, toUID))
			}
		}
		if fromUser == header.Uname {
			header.Uname = toUser
			if verbose {
				effectiveChanges++
				verboseOut = append(verboseOut, fmt.Sprintf("user:%s:%s", fromUser, toUser))
			}
		}
		if fromGID == header.Gid {
			header.Gid = toGID
			if verbose {
				effectiveChanges++
				verboseOut = append(verboseOut, fmt.Sprintf("GID:%d:%d", fromGID, toGID))
			}
		}
		if fromGroup == header.Gname {
			header.Gname = toGroup
			if verbose {
				effectiveChanges++
				verboseOut = append(verboseOut, fmt.Sprintf("group:%s:%s", fromGroup, toGroup))
			}
		}
		if err := tarOutput.WriteHeader(header); err != nil {
			panic(err)
		}
		io.CopyN(tarOutput, tarInput, header.Size)
		if verbose {
			if len(verboseOut) > 0 {
				entryChanges++
				log.Println(verboseOut, header.Name)
			}
		}
	}
	if err := tarOutput.Flush(); err != nil {
		panic(err)
	}
	if err := tarOutput.Close(); err != nil {
		panic(err)
	}
	if verbose {
		log.Printf("Changed %d datapoints amongst %d/%d entries", effectiveChanges, entryChanges, totalEntries)
	}
}
