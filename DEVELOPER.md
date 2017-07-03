### Default values for options

Go forces default values for types like bool, int and strings.

### Expected breaking changes while developing v4.0

* Timeouts and Deadlines currently implemented in a way they done in Java/PHP
SDKs. The problem is that there is no option to reassign certain values during
existing client session without recreating it. So http.Client and http.Transport
initialization and interaction should be completely reviewed.
