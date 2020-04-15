# Goldilox

Utility to roughly identify max header size via GET on a resource.  Conducts a 
binary search across the provided range to reduce invocations on the specified 
resource.

## Usage

```bash
Usage of ./goldilox:
  -max int
        maximum header size (bytes) for tesing range (defaults to 10000) (default 10000)
  -min int
        minimum header size (bytes) for tesing range (defaults to 1) (default 1)
  -url string
        url to send GET requests to (required)
```
