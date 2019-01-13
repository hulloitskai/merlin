package info

var (
	// Version is the program version, set during compile time using:
	//   -ldflags -X github.com/stevenxie/merlin/api/internal.Version=$(VERSION)
	Version = "unset"

	// Namespace is the project namespace, to be used as prefixes for environment
	// variables, etc.
	Namespace = "merlin"
)
