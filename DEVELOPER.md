### Default values for options

Go forces default values for types like bool, int and strings.

### Some breaking changes from the previous v4's
	⁃	config field UseHttp2 is now UseHTTP2
	⁃	config field Uuid is now UUID
	⁃	Get State/ WhereNow Uuid is now UUID
	⁃	In Fire/Publish Ttl() is now TTL()
	⁃	In Grant Ttl() is now TTL()
	⁃	PNPAMEntityData Ttl is now TTL
	⁃	PNAccessManagerKeyData Ttl is now TTL
	⁃	TlsEnabled is now TLSEnabled in StatusResponse and ResponseInfo
	⁃	Uuid is now UUID in StatusResponse and ResponseInfo