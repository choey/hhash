# hhash

hhash is a consistent, human-readable hashing engine

# why
author desired something like the 
[docker container name generator](https://github.com/moby/moby/blob/master/pkg/namesgenerator/names-generator.go)
but with consistent hashing and rich custom pattern 

# installation
```
$ go get github.com/choey/hhash
$ go get github.com/choey/hhash/cmd/hhash
```

# features
* supports custom pattern
* generates random human-readable string
* generates human-readable hash
* reports collision rate for the given datasets

# futures
* add dataset for person names
* scrub datasets to be less offensive

# examples
* random string (which is just the hash of a random UUID)
```
$ hhash
cedarn_kappa
```
* custom pattern (upper case adverb followed by uppercase adjective)
```
$ hhash -p '[%n-%n]'
[devon-freeware]
```
* value hash
```
$ hhash -s 'hello, world' -v
2019/02/07 23:36:09 using pattern: %j_%n
2019/02/07 23:36:09 elapsed time 373.439µs
fond_bonefish
$ hhash -s 'hello, world' -v
2019/02/07 23:36:33 using pattern: %j_%n
2019/02/07 23:36:33 elapsed time 331.244µs
fond_bonefish
```
(verbose mode calculates collision rate and thus is slightly slower than non-verbose mode)

# attributions
* [WordNet](https://wordnet.princeton.edu/license-and-commercial-use) provides the default datasets
