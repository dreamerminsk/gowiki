package atomic

// nocmp is an uncomparable struct. Embed this inside another struct to make
// it uncomparable.
//
//  type Foo struct {
//    nocmp
//    // ...
//  }
//
// This DOES NOT:
//
//  - Disallow shallow copies of structs
//  - Disallow comparison of pointers to uncomparable structs
type nocmp [0]func()
