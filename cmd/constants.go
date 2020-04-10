package cmd

const (
	NXTalkProxyDDefaultLaddr = ":1969"                      // NXTalkProxyDDefaultLaddr is the default Host:port of `nxtalkproxyd`.
	ConfigurationFileDocs    = "Configuration file to use." // ConfigurationFileDocs is the documentation for the configuration file flag.
)

const (
	CouldNotBindFlagsErrorMessage        = "could not bind flags"         // CouldNotBindFlagsErrorMessage is the error message to throw if binding the flags has failed.
	CouldNotStartRootCommandErrorMessage = "could not start root command" // CouldNotStartRootCommandErrorMessage is the error message to throw if starting the root command has failed.
)
