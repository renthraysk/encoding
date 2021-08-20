Encoding

Parsing HTTP Accept-Encoding headers

Points to note:

- This code only cares whether a useragent supports a compression method or not. So q values of 0 means explictly unsupported, and positive q values means supported.
- Parse() is allocation free.
