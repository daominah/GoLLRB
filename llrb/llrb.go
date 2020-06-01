// Copyright 2010 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A Left-Leaning Red-Black (LLRB) implementation of 2-3 balanced binary search trees,
// based on the following work:
//
//   http://www.cs.princeton.edu/~rs/talks/LLRB/08Penn.pdf
//   http://www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf
//   http://www.cs.princeton.edu/~rs/talks/LLRB/Java/RedBlackBST.java
//
//  2-3 trees (and the run-time equivalent 2-3-4 trees) are the de facto standard BST
//  algoritms found in implementations of Python, Java, and other libraries. The LLRB
//  implementation of 2-3 trees is a recent improvement on the traditional implementation,
//  observed and documented by Robert Sedgewick.
//
package llrb

import (
	"fmt"
	"strings"
)

// LLRB is an order statistic tree,
// this is an augmented Left-Leaning Red-Black (LLRB) implementation of 2-3 trees
type LLRB struct {
	count int
	root  *Node
}

type Node struct {
	Item
	Left, Right *Node // Pointers to left and right child nodes
	Black       bool  // If set, the color of the link (incoming from the parent) is black
	// In the LLRB, new nodes are always red, hence the zero-value for node

	// size of the subtree that has root is this Node,
	// NDescendants == tree_count in for the tree's root Node
	NDescendants int
}

type Item interface {
	Less(than Item) bool
}

//
func less(x, y Item) bool {
	if x == pinf {
		return false
	}
	if x == ninf {
		return true
	}
	return x.Less(y)
}

// Inf returns an Item that is "bigger than" any other item, if sign is positive.
// Otherwise  it returns an Item that is "smaller than" any other item.
func Inf(sign int) Item {
	if sign == 0 {
		panic("sign")
	}
	if sign > 0 {
		return pinf
	}
	return ninf
}

var (
	ninf = nInf{}
	pinf = pInf{}
)

type nInf struct{}

func (nInf) Less(Item) bool {
	return true
}

type pInf struct{}

func (pInf) Less(Item) bool {
	return false
}

// New allocates a new tree
func New() *LLRB {
	return &LLRB{}
}

// SetRoot sets the root node of the tree.
// It is intended to be used by functions that deserialize the tree.
func (t *LLRB) SetRoot(r *Node) {
	t.root = r
}

// Root returns the root node of the tree.
// It is intended to be used by functions that serialize the tree.
func (t *LLRB) Root() *Node {
	return t.root
}

// Len returns the number of nodes in the tree.
func (t *LLRB) Len() int { return t.count }

// Has returns true if the tree contains an element whose order is the same as that of key.
func (t *LLRB) Has(key Item) bool {
	return t.Get(key) != nil
}

// Get retrieves an element from the tree whose order is the same as that of key.
func (t *LLRB) Get(key Item) Item {
	h := t.root
	for h != nil {
		switch {
		case less(key, h.Item):
			h = h.Left
		case less(h.Item, key):
			h = h.Right
		default:
			return h.Item
		}
	}
	return nil
}

// Min returns the minimum element in the tree.
func (t *LLRB) Min() Item {
	h := t.root
	if h == nil {
		return nil
	}
	for h.Left != nil {
		h = h.Left
	}
	return h.Item
}

// Max returns the maximum element in the tree.
func (t *LLRB) Max() Item {
	h := t.root
	if h == nil {
		return nil
	}
	for h.Right != nil {
		h = h.Right
	}
	return h.Item
}

func (t *LLRB) ReplaceOrInsertBulk(items ...Item) {
	for _, i := range items {
		t.ReplaceOrInsert(i)
	}
}

func (t *LLRB) InsertNoReplaceBulk(items ...Item) {
	for _, i := range items {
		t.InsertNoReplace(i)
	}
}

// ReplaceOrInsert inserts item into the tree. If an existing
// element has the same order, it is removed from the tree and returned.
func (t *LLRB) ReplaceOrInsert(item Item) Item {
	// TODO: correct NDescendants
	if item == nil {
		panic("inserting nil item")
	}
	var replaced Item
	t.root, replaced = t.replaceOrInsert(t.root, item)
	t.root.Black = true
	if replaced == nil {
		t.count++
	}
	return replaced
}

func (t *LLRB) replaceOrInsert(h *Node, item Item) (*Node, Item) {
	if h == nil {
		return newNode(item), nil
	}

	h = walkDownRot23(h)

	var replaced Item
	if less(item, h.Item) { // BUG
		h.Left, replaced = t.replaceOrInsert(h.Left, item)
	} else if less(h.Item, item) {
		h.Right, replaced = t.replaceOrInsert(h.Right, item)
	} else {
		replaced, h.Item = h.Item, item
	}

	h = walkUpRot23(h)

	return h, replaced
}

// InsertNoReplace inserts item into the tree. If an existing
// element has the same order, both elements remain in the tree.
func (t *LLRB) InsertNoReplace(item Item) {
	if item == nil {
		panic("inserting nil item")
	}
	t.root = t.insertNoReplace(t.root, item)
	t.root.Black = true
	t.count++
}

func (t *LLRB) insertNoReplace(h *Node, item Item) *Node {
	if h == nil {
		return newNode(item)
	}

	h = walkDownRot23(h)

	h.NDescendants += 1
	if less(item, h.Item) {
		h.Left = t.insertNoReplace(h.Left, item)
	} else {
		h.Right = t.insertNoReplace(h.Right, item)
	}

	return walkUpRot23(h)
}

// Rotation driver routines for 2-3 algorithm

func walkDownRot23(h *Node) *Node { return h }

func walkUpRot23(h *Node) *Node {
	if isRed(h.Right) && !isRed(h.Left) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}

// Rotation driver routines for 2-3-4 algorithm

func walkDownRot234(h *Node) *Node {
	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}

func walkUpRot234(h *Node) *Node {
	if isRed(h.Right) && !isRed(h.Left) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	return h
}

// DeleteMin deletes the minimum element in the tree and returns the
// deleted item or nil otherwise.
func (t *LLRB) DeleteMin() Item {
	// TODO: correct NDescendants
	var deleted Item
	t.root, deleted = deleteMin(t.root)
	if t.root != nil {
		t.root.Black = true
	}
	if deleted != nil {
		t.count--
	}
	return deleted
}

// deleteMin code for LLRB 2-3 trees
func deleteMin(h *Node) (*Node, Item) {
	if h == nil {
		return nil, nil
	}
	if h.Left == nil {
		return nil, h.Item
	}

	if !isRed(h.Left) && !isRed(h.Left.Left) {
		h = moveRedLeft(h)
	}

	var deleted Item
	h.Left, deleted = deleteMin(h.Left)

	return fixUp(h), deleted
}

// DeleteMax deletes the maximum element in the tree and returns
// the deleted item or nil otherwise
func (t *LLRB) DeleteMax() Item {
	// TODO: correct NDescendants
	var deleted Item
	t.root, deleted = deleteMax(t.root)
	if t.root != nil {
		t.root.Black = true
	}
	if deleted != nil {
		t.count--
	}
	return deleted
}

func deleteMax(h *Node) (*Node, Item) {
	if h == nil {
		return nil, nil
	}
	if isRed(h.Left) {
		h = rotateRight(h)
	}
	if h.Right == nil {
		return nil, h.Item
	}
	if !isRed(h.Right) && !isRed(h.Right.Left) {
		h = moveRedRight(h)
	}
	var deleted Item
	h.Right, deleted = deleteMax(h.Right)

	return fixUp(h), deleted
}

// Delete deletes an item from the tree whose key equals key.
// The deleted item is return, otherwise nil is returned.
func (t *LLRB) Delete(key Item) Item {
	// TODO: correct NDescendants
	var deleted Item
	t.root, deleted = t.delete(t.root, key)
	if t.root != nil {
		t.root.Black = true
	}
	if deleted != nil {
		t.count--
	}
	return deleted
}

func (t *LLRB) delete(h *Node, item Item) (*Node, Item) {
	var deleted Item
	if h == nil {
		return nil, nil
	}
	if less(item, h.Item) {
		if h.Left == nil { // item not present. Nothing to delete
			return h, nil
		}
		if !isRed(h.Left) && !isRed(h.Left.Left) {
			h = moveRedLeft(h)
		}
		h.Left, deleted = t.delete(h.Left, item)
	} else {
		if isRed(h.Left) {
			h = rotateRight(h)
		}
		// If @item equals @h.Item and no right children at @h
		if !less(h.Item, item) && h.Right == nil {
			return nil, h.Item
		}
		// PETAR: Added 'h.Right != nil' below
		if h.Right != nil && !isRed(h.Right) && !isRed(h.Right.Left) {
			h = moveRedRight(h)
		}
		// If @item equals @h.Item, and (from above) 'h.Right != nil'
		if !less(h.Item, item) {
			var subDeleted Item
			h.Right, subDeleted = deleteMin(h.Right)
			if subDeleted == nil {
				panic("logic")
			}
			deleted, h.Item = h.Item, subDeleted
		} else { // Else, @item is bigger than @h.Item
			h.Right, deleted = t.delete(h.Right, item)
		}
	}

	return fixUp(h), deleted
}

// Internal node manipulation routines

func newNode(item Item) *Node {
	return &Node{
		Item:         item,
		NDescendants: 1,
	}
}

func isRed(h *Node) bool {
	if h == nil {
		return false
	}
	return !h.Black
}

func rotateLeft(h *Node) *Node {
	parentSize := h.NDescendants
	leftChildSize := size(h.Left)
	rightChildL1LeftChildL2Size := size(h.Right.Left)

	x := h.Right
	if x.Black {
		panic("rotating a black link")
	}
	h.Right = x.Left
	x.Left = h
	x.Black = h.Black
	h.Black = false

	x.NDescendants = parentSize
	h.NDescendants = leftChildSize + rightChildL1LeftChildL2Size + 1

	return x
}

func rotateRight(h *Node) *Node {
	parentSize := h.NDescendants
	rightChildSize := size(h.Right)
	leftChildL1rightChildL2Size := size(h.Left.Right)

	x := h.Left
	if x.Black {
		panic("rotating a black link")
	}
	h.Left = x.Right
	x.Right = h
	x.Black = h.Black
	h.Black = false

	x.NDescendants = parentSize
	h.NDescendants = rightChildSize + leftChildL1rightChildL2Size

	return x
}

// flip changes color of the node and its children,
// REQUIRE: Left and Right children must be present
func flip(h *Node) {
	h.Black = !h.Black
	h.Left.Black = !h.Left.Black
	h.Right.Black = !h.Right.Black
}

// REQUIRE: Left and Right children must be present
func moveRedLeft(h *Node) *Node {
	flip(h)
	if isRed(h.Right.Left) {
		h.Right = rotateRight(h.Right)
		h = rotateLeft(h)
		flip(h)
	}
	return h
}

// REQUIRE: Left and Right children must be present
func moveRedRight(h *Node) *Node {
	flip(h)
	if isRed(h.Left.Left) {
		h = rotateRight(h)
		flip(h)
	}
	return h
}

func fixUp(h *Node) *Node {
	if isRed(h.Right) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}

// size is convenient to get node_NDescendants (node can be nil)
func size(h *Node) int {
	if h == nil {
		return 0
	}
	return h.NDescendants
}

func (h *Node) String() string {
	if h != nil {
		return fmt.Sprintf("[k:%v,%v,%v]",
			h.Item, h.NDescendants, h.Black)
	} else {
		return "nil"
	}
}

func (t *LLRB) printBFS() string {
	lines := make([]string, 0)
	visiteds := make(map[*Node]bool, t.count)
	type QueueElem struct {
		node   *Node
		parent string
	}
	q := []QueueElem{{node: t.root, parent: "IAmRoot"}}
	for len(q) > 0 {
		pop := q[0]
		q = q[1:]
		visiteds[pop.node] = true
		parentStr := fmt.Sprintf("%v", pop.node.Item)
		if pop.node.Left != nil {
			q = append(q, QueueElem{node: pop.node.Left, parent: parentStr})
		}
		if pop.node.Right != nil {
			q = append(q, QueueElem{node: pop.node.Right, parent: parentStr})
		}
		line := fmt.Sprintf("parent: %v, node: %v, ", pop.parent, pop.node)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// GetByRank retrieves an Item with a given rank r (rank start from 1).
// this func only returns nil if the tree has length 0 or the tree is invalid.
func (t *LLRB) GetByRank(r int) Item {
	node := t.getByRank(t.root, r)
	if node == nil {
		if r <= 0 {
			return t.Min()
		} else { // r > tree_length
			return t.Max()
		}
	}
	return node.Item
}

func (t *LLRB) getByRank(h *Node, r int) *Node {
	if h == nil {
		return nil
	}
	hRank := size(h.Left) + 1
	if r == hRank {
		return h
	}
	if r < hRank {
		if h.Left == nil { // never expected to reach this branch
			return nil
		}
		return t.getByRank(h.Left, r)
	}
	if h.Right == nil { // never expected to reach this branch
		return nil
	}
	return t.getByRank(h.Right, r-hRank)
}

// GetRankOf determines rank of an key (rank start from 1),
// this func returns the rank and one Item in the tree that equal to key
func (t *LLRB) GetRankOf(key Item) (int, Item) {
	path := t.get(key)
	//fmt.Println("path: ", path)
	return t.getRankOf(path, key)
}

// get returns path from the root to a node whose order is the same as that of key,
// path[last] is nil if the tree does not contain exact key Item
func (t *LLRB) get(key Item) []*Node {
	path := make([]*Node, 0)
	h := t.root
	for h != nil {
		path = append(path, h)
		switch {
		case less(key, h.Item):
			h = h.Left
		case less(h.Item, key):
			h = h.Right
		default: // exactly equal
			return path
		}
	}
	path = append(path, nil)
	return path
}

func (t *LLRB) getRankOf(path []*Node, key Item) (int, Item) {
	if len(path) < 1 {
		return 0, nil
	}
	foundItem := path[len(path)-1]
	var r int
	if foundItem != nil {
		r = size(foundItem.Left) + 1
	} else { // return rank of nearest parent node for non-existed item
		if len(path) < 2 {
			return 0, nil
		}
		r = size(path[len(path)-2].Left) + 1
	}
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == nil {
			continue
		}
		if i-1 < 0 || path[i-1] == nil {
			continue
		}
		if path[i] == path[i-1].Right { // if current node is a right node
			r += size(path[i-1].Left) + 1 // add size of the left sibling to the rank
		}
	}
	return r, foundItem
}
