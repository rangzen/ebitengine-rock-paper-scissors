package main

import "errors"

// Beware,
// the following code search the nearest neighbor but from a different type
// than object given in parameter (index).

// Circular implements an index neighbor search.
// It is for testing purpose only.
type Circular struct{}

func (l *Circular) Neighbor(objs []Obj, i int) (int, error) {
	for j := 1; j < len(objs)-1; j++ {
		n := (i + j) % len(objs)
		if objs[i].t != objs[n].t {
			return n, nil
		}
	}
	return 0, errors.New("no neighbor found")
}

// Linear implements a linear search of the nearest neighbor.
// https://en.wikipedia.org/wiki/Nearest_neighbor_search#Linear_search
type Linear struct{}

func (l *Linear) Neighbor(objs []Obj, i int) (int, error) {
	bestIndex := -1
	bestDist := screenWidth * screenWidth * screenHeight * screenHeight
	for j := 1; j < len(objs)-1; j++ {
		n := (i + j) % len(objs)
		if objs[i].t != objs[n].t {
			dist := sqrDist(objs[i], objs[n])
			if dist < bestDist {
				bestIndex = n
				bestDist = dist
			}
		}
	}
	if bestIndex == -1 {
		return 0, errors.New("no neighbor found")
	}
	return bestIndex, nil
}

// sqrDist returns the square of the distance between two objects.
// To avoid a square root, we compare the square of the distance.
func sqrDist(o1 Obj, o2 Obj) int {
	return (o1.x-o2.x)*(o1.x-o2.x) + (o1.y-o2.y)*(o1.y-o2.y)
}
