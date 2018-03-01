# lwwset
LWW-Element-Set CRDT implementation with Go.

## Install
To use lwwset, first:

```Go
go get github.com/drkaka/lwwset
```

## Usage

Create a set:

```Go
set := lwwset.NewSet()
```

Add an element with timestamp

```Go
set.Add(element, ts)
```

Lookup an element **(when timestamps in both addSet & removeSet are equal, the element exists)**

```Go
exist := set.Lookup(element)
```

Remove an element at the timestamp

```Go
set.Remove(element, ts)
```

Merge another set
```Go
set.Merge(set2)
```
