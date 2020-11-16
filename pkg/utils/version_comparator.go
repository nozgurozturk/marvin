package utils

import "fmt"

// Compares semantic version of packages
// 1.0.0 -> major.minor.patch
// If it's true, package is outdated
func CompareVersions(latestVersion string, currentVersion string) bool {
	var lMajor, lMinor, lPatch int
	var cMajor, cMinor, cPatch int
	// Parses registry version numbers to integer
	fmt.Sscanf(latestVersion, "%d.%d.%d", &lMajor, &lMinor, &lPatch)
	// Parses repository version numbers to integer
	fmt.Sscanf(currentVersion, "%d.%d.%d", &cMajor, &cMinor, &cPatch)


	if cMajor < lMajor { 	// Check major
		return true
	} else if cMinor < lMinor { 	// Check minor
		return true
	} else if cPatch < lPatch { 	// Check patch
		return true
	}
	return false
}
