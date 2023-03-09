package pool

const (
	writeGeolocEntryQuery = `
		INSERT INTO geoloc.entry(ip, code, country, city, lat, lon)
		VALUES (?, ?, ?, ?, ?, ?)`

	readGeolocEntryQuery = `
		SELECT code, country, city, lat, lon
		FROM geoloc.entry
		WHERE ip = ?`
)
