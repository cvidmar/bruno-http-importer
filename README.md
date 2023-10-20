# Bruno Http Importer

A migration tool to translate a tree of directories containing .http files into one with .bru files for Bruno. Bruno is an "Opensource IDE for exploring and testing APIs". See https://github.com/usebruno/bruno

By .http files we mean files like:

```
POST https://www.example.com/one
Cookie: session=1234567890
Content-Type: application/json

{
"foo": "bar"
}
```

Current limitations:

- Only migrates POST and GET methods
- POST requests migrations only support JSON body
