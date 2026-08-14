//go:build !linux

package linux

func filesystemUsage(string) (used, total int64) { return 0, 0 }
