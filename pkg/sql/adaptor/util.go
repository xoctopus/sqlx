package adaptor

import "net/url"

func DatabaseNameFromDSN(u *url.URL) string {
	database := u.Path
	if len(database) > 0 && database[0] == '/' {
		database = database[1:]
	}
	return database
}
