package requests

import (
	"path/filepath"
	"strings"
)

// Default MIME types for unknown text or
// binary files.
//
// If the MIME type can't be guessed from the MIMETypes map,
// use one of these defaults, depending on if the file is
// text or binary.
const (
	MIMEDefaultText   = "text/plain"
	MIMEDefaultBinary = "application/octet-stream"
)

// MIMETypes maps file extensions to MIME types.
// The map keys are lower-cased and include the leading period.
//
// The data comes from the MDN Web Docs page for common MIME types:
//
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
//
// The data isn't meant to be complete but to be a good starting point
// for helping the user guess an unknown MIME type.
//
// There are somes caveats listed in the MDN page:
//   - ".mid" is set to "audio/midi" but could also be "audio/x-midi"
//   - ".3gp" is set to "video/3gpp" but could also be "audio/3gpp" if it doesn't contain video
//   - ".3g2" is set to "video/3gpp2" but could also be "audio/3gpp2" if it doesn't contain video
//   - ".xml" is set to "application/xml" but in some older implementations could also be "text/xml"
//
// If you can't find a sutable MIME type for your file from this map, try
// using either of the supplied defaults, depending on if the file is binary or not:
//
// 		const (
//			MIMEDefaultText   = "text/plain"
//			MIMEDefaultBinary = "application/octet-stream"
// 		)
//
var MIMETypes = map[string]string{
	".aac":    "audio/aac",
	".abw":    "application/x-abiword",
	".arc":    "application/x-freearc",
	".avi":    "video/x-msvideo",
	".azw":    "application/vnd.amazon.ebook",
	".bin":    "application/octet-stream",
	".bmp":    "image/bmp",
	".bz":     "application/x-bzip",
	".bz2":    "application/x-bzip2",
	".cda":    "application/x-cdf",
	".csh":    "application/x-csh",
	".css":    "text/css",
	".csv":    "text/csv",
	".doc":    "application/msword",
	".docx":   "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".eot":    "application/vnd.ms-fontobject",
	".epub":   "application/epub+zip",
	".gz":     "application/gzip",
	".gif":    "image/gif",
	".htm":    "text/html",
	".html":   "text/html",
	".ico":    "image/vnd.microsoft.icon",
	".ics":    "text/calendar",
	".jar":    "application/java-archive",
	".jpeg":   "image/jpeg",
	".jpg":    "image/jpeg",
	".js":     "text/javascript", // (Specifications: HTML and its reasoning, and IETF)
	".json":   "application/json",
	".jsonld": "application/ld+json",
	".mid":    "audio/midi", // could also be "audio/x-midi"
	".midi":   "audio/midi",
	".mjs":    "text/javascript",
	".mp3":    "audio/mpeg",
	".mp4":    "video/mp4",
	".mpeg":   "video/mpeg",
	".mpkg":   "application/vnd.apple.installer+xml",
	".odp":    "application/vnd.oasis.opendocument.presentation",
	".ods":    "application/vnd.oasis.opendocument.spreadsheet",
	".odt":    "application/vnd.oasis.opendocument.text",
	".oga":    "audio/ogg",
	".ogv":    "video/ogg",
	".ogx":    "application/ogg",
	".opus":   "audio/opus",
	".otf":    "font/otf",
	".png":    "image/png",
	".pdf":    "application/pdf",
	".php":    "application/x-httpd-php",
	".ppt":    "application/vnd.ms-powerpoint",
	".pptx":   "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	".rar":    "application/vnd.rar",
	".rtf":    "application/rtf",
	".sh":     "application/x-sh",
	".svg":    "image/svg+xml",
	".swf":    "application/x-shockwave-flash",
	".tar":    "application/x-tar",
	".tif":    "image/tiff",
	".tiff":   "image/tiff",
	".ts":     "video/mp2t",
	".ttf":    "font/ttf",
	".txt":    "text/plain",
	".vsd":    "application/vnd.visio",
	".wav":    "audio/wav",
	".weba":   "audio/webm",
	".webm":   "video/webm",
	".webp":   "image/webp",
	".woff":   "font/woff",
	".woff2":  "font/woff2",
	".xhtml":  "application/xhtml+xml",
	".xls":    "application/vnd.ms-excel",
	".xlsx":   "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".xml":    "application/xml", // is recommended as of RFC 7303 (section 4.1), but text/xml is still seen sometimes. A file with the extension .xml can often be given a more specific MIME type depending on how its contents are meant to be interpreted (for instance, an Atom feed is application/atom+xml), but application/xml serves as a valid default.
	".xul":    "application/vnd.mozilla.xul+xml",
	".zip":    "application/zip",
	".3gp":    "video/3gpp",  // or "audio/3gpp" if it doesn't contain video
	".3g2":    "video/3gpp2", // or "audio/3gpp2" if it doesn't contain video
	".7z":     "application/x-7z-compressed",
}

// GuessMIMEType is a helper function for guessing the MIME type
// of a file, based on it's filename using the MIMETypes map.
//
// It returns the MIME type and a boolean indicating if it was
// in the map or not.
//
// If a MIME type can't be determined from the filename, consider
// using either MIMEDefaultText or MIMEDefaultBinary, depending
// on if the file is binary or not.
func GuessMIME(filename string) (string, bool) {
	ext := strings.ToLower(filepath.Ext(filename))
	mime, ok := MIMETypes[ext]
	return mime, ok
}
