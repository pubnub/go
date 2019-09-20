Concise Binary Object Representation (CBOR) is a superset of JSON's schema that's faster and more compact.

* http://tools.ietf.org/html/rfc7049
* http://cbor.io/

import cbor "github.com/brianolson/cbor_go"

This implementation attempts to do serialization to/from struct types using reflection, but doesn't do 100% of cases like that right. It _should_ do everything fine serializing `map[string]interface{}` and `[]interface{}` and other basic types. It passes the test of decoding 100% of CBOR common appendix test strings.

---

For Python implementation see:
https://github.com/brianolson/cbor_py