package requestConfig

const (
	BlockStart = "---Request"

	CommentLineStartsWith = "#"

	SequenceLineStartsWith                = ":"
	SequenceLineToUpperSequenceStartMatch = "SEQUENCE"

	ExpectLineStartsWith      = "expect:"
	ExpectStatusCodeSeparator = ","

	MethodAndUriLineMethodRestSeparator = " "
	UriStart                            = "/"

	HeaderLineNameValueSeparator = ":"

	AuthorizationHeaderName = "Authorization"
)
