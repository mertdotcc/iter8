package core

// lookup.go -(internal) implementation of gRPC Lookup method

import (
	"errors"
	"fmt"
	"hash/maphash"

	abnapp "github.com/iter8-tools/iter8/abn/application"
	"github.com/iter8-tools/iter8/base/log"
)

var versionHasher maphash.Hash

// lookupInternal is detailed implementation of gRPC method Lookup
func lookupInternal(application string, user string) (*abnapp.Application, *string, error) {
	// if user is not provided, fail
	if user == "" {
		return nil, nil, errors.New("no user session provided")
	}

	// check that we have a record of the application
	a, err := abnapp.Applications.Get(application)
	if err != nil || a == nil {
		return nil, nil, fmt.Errorf("application %s not found", application)
	}

	// use rendezvous hash to get track for user, fail if not present
	abnapp.Applications.RLock(application)
	defer abnapp.Applications.RUnlock(application)
	track := rendezvousGet(a, user)
	return a, &track, nil
}

// rendezvousGet is an implementation of rendezvous hashing (cf. https://en.wikipedia.org/wiki/Rendezvous_hashing)
// It returns a consistent track for a given application and user combination.
// The track is chosen uniformly at random from among the current set of tracks
// associated with an application.
// We want to always return the same track for the same user so long as the
// application remains unchanged -- there are no change in the set of versions
// and no change to the track mapping.
// To do this, we hash the combination of user and version. We don't use the track identifier
// because the track identifier is associated with multiple versions over time; we do not
// require a fixed mapping when this mapping changes.
// We select the version, user pair with the largest hash value ("score").
// Inspired by https://github.com/tysonmote/rendezvous/blob/master/rendezvous.go
func rendezvousGet(a *abnapp.Application, user string) string {
	// current maximimum score as computed by the hash function
	var maxScore uint64
	// maxTrack is the track with the current maximum score
	var maxTrack string
	// maxVersion is the version name associated with maxTrack
	var maxVersion string

	for track, version := range a.Tracks {
		score := hash(version, user)
		log.Logger.Debugf("hash(%s,%s) --> %d  --  %d", version, user, score, maxScore)
		if score > maxScore || (score == maxScore && version > maxVersion) {
			maxScore = score
			maxVersion = version
			maxTrack = track
		}
	}
	return maxTrack
}

// hash computes the score for a version, user combination
func hash(version, user string) uint64 {
	versionHasher.Reset()
	_, _ = versionHasher.WriteString(user)
	_, _ = versionHasher.WriteString(version)
	return versionHasher.Sum64()
}
