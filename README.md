# Goldilox

Utility to roughly identify max header size via GET on a resource.  Conducts a 
binary search across the provided range to reduce invocations on the specified 
resource.

```
Usage of ./goldilox:
  -chunkSize int
        size (bytes) to partition header into separate key/value paris (default 3000)
  -max int
        maximum header size (bytes) for tesing range (default 10000)
  -min int
        minimum header size (bytes) for tesing range (default 1)
  -url string
        url to send GET requests to (required)
```
