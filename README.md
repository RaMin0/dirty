# Dirty

Track changes to struct fields.

## Example

```go
u := User{}


dirty.Track(&u)

dirty.Changed(&u) // => false


u.Name = "John Doe"


dirty.Changed(&u) // => true

dirty.Changes(&u)
// => map[string][]interface{} {
//      "Name" => ["", "John Doe"]
//    }


dirty.Forget(&u)

dirty.Changed(&u) // => panic!
```
