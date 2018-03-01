# lwwset [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov]
LWW-Element-Set CRDT thread-safe implementation with Go.

## Install
To use lwwset, first:

```Go
go get github.com/drkaka/lwwset
```

## Usage

**Init** a set:

```Go
set := lwwset.NewSet()
```

**Add** an element with timestamp

```Go
set.Add(element, ts)
```

**Lookup** an element **(when timestamps in addSet & removeSet are equal, the element exists)**

```Go
exist := set.Lookup(element)
```

**Remove** an element at the timestamp

```Go
set.Remove(element, ts)
```

**Merge** another set
```Go
set.Merge(set2)
```


[ci-img]: https://travis-ci.org/drkaka/lwwset.svg?branch=master
[ci]: https://travis-ci.org/drkaka/lwwset
[cov-img]: https://coveralls.io/repos/github/drkaka/lwwset/badge.svg?branch=master
[cov]: https://coveralls.io/github/drkaka/lwwset?branch=master