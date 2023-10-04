package global

// Build ldflags.
var (
	// VERS is the latest cryptor version tag. Set by linker -ldflags "-X main.VERS=..."
	VERS = "v0.2.0"
	// OS is the target operating system and architecture. Set by linker -ldflags "-X main.OS=..."
	OS = "-"
	// BUILD is the date the executable was built.
	BUILT = "-"
	// COMMIT is the Git commit hash.
	COMMIT = "-"

	// URL to fetch exchange rates.
	// XRATES_QUERY = "https://api.exchangerate.host/latest?base=usd&access_key=YOUR_ACCESS_KEY"
	XRATES_QUERY = "https://openexchangerates.org/api/latest.json?app_id=404d2ec9a36a4f73948dccb71887b788"
)
