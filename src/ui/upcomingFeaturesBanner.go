package ui

import (
	"fmt"
	"strings"

	"servercommander/src/utils"
)

func UpcomingFeaturesBanner() {
	features := map[string]string{
		"Utility": "Various utility commands needed for a good console: htop, clear, help, log, etc.",
		"SSH":     "Remote connection to another server: SSH remote connection.",
		"FTP":     "Transfer data over a remote connection: FTP/sFTP.",
	}

	maxFeatureLength := 0
	maxDescriptionLength := 0

	for feature, description := range features {
		if len(feature) > maxFeatureLength {
			maxFeatureLength = len(feature)
		}
		if len(description) > maxDescriptionLength {
			maxDescriptionLength = len(description)
		}
	}

	fmt.Println(utils.Cyan, "Upcoming Features:", utils.Reset)
	separator := fmt.Sprintf("+-%s-+-%s-+", strings.Repeat("-", maxFeatureLength), strings.Repeat("-", maxDescriptionLength))
	fmt.Println(separator)
	fmt.Printf("| %-*s | %-*s |\n", maxFeatureLength, "Feature", maxDescriptionLength, "Description")
	fmt.Println(separator)

	for feature, description := range features {
		fmt.Printf("| %-*s | %-*s |\n", maxFeatureLength, feature, maxDescriptionLength, description)
	}

	fmt.Println(separator)
}
