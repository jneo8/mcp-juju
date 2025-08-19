package jujuadapter

import (
	"fmt"
	"strconv"
	"strings"
)

// parseURIParameters extracts parameters from a URI based on the resource template configuration
func parseURIParameters(uri string, config ResourceTemplateConfig) (map[string]string, error) {
	params := make(map[string]string)

	// Remove protocol and split into parts
	uriPath := strings.TrimPrefix(uri, "juju://")
	parts := strings.Split(uriPath, "/")

	// Use the resource template configuration to map URI parts to parameters
	// The config tells us exactly which URI parts map to which parameters
	for paramName, argIndexStr := range config.URIToArgs {
		argIndex, err := strconv.Atoi(argIndexStr)
		if err != nil {
			return nil, fmt.Errorf("invalid argument index '%s' for parameter '%s'", argIndexStr, paramName)
		}
		
		// Adjust index to account for the base path (e.g., "config" is index 0)
		uriIndex := argIndex + 1 // +1 because index 0 is "config"
		
		if uriIndex < len(parts) && parts[uriIndex] != "" {
			params[paramName] = parts[uriIndex]
		}
		// If the parameter is not present, that's OK for optional parameters
	}

	// Also handle any URI-to-flag mappings
	for paramName, flagName := range config.URIToFlags {
		// For flags, we would extract from URI parts as well
		// This is extensible for future templates that map URI parts to flags
		// For now, this is just a placeholder for future functionality
		_ = paramName
		_ = flagName
	}

	return params, nil
}