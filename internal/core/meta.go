package core

import (
	"runtime"
	"runtime/debug"
	"time"
)

// Meta contains application metadata (version, os, build info, etc).
type Meta struct {
	BuildCommit    string
	BuildTime      time.Time
	BuildGoVersion string
	BuildVersion   string
	BuildChecksum  string
	Version        string
	GOOS           string
	GOARCH         string
}

// NewMeta returns a new Meta struct.
func NewMeta(version, commit, date string) *Meta {
	buildTime, _ := time.Parse(time.RFC3339, date)

	meta := &Meta{
		BuildCommit: commit,
		BuildTime:   buildTime,
		Version:     version,
		GOOS:        runtime.GOOS,
		GOARCH:      runtime.GOARCH,
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		meta.BuildGoVersion = info.GoVersion
		meta.BuildVersion = info.Main.Version
		meta.BuildChecksum = info.Main.Sum
	}

	return meta
}
