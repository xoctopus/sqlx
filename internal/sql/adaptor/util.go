package adaptor

import "net/url"

// DatabaseNameFromDSN returns the database name from u.Path, stripping a leading '/'.
func DatabaseNameFromDSN(u *url.URL) string {
	database := u.Path
	if len(database) > 0 && database[0] == '/' {
		database = database[1:]
	}
	return database
}
