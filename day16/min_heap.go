package main

import "errors"

type MinHeap struct {
	data   []int
	length int
}

func (mh *MinHeap) Insert(num int) {
	if mh.length < len(mh.data) {
		mh.data[mh.length] = num
	} else {
		mh.data = append(mh.data, num)
	}
	mh.heapifyUp(mh.length)
	mh.length += 1
}

func (mh *MinHeap) Delete() (int, error) {

	if mh.length == 0 {
		return 0, errors.New("Empty")
	}

	out := mh.data[0]
	mh.length -= 1

	if mh.length == 0 {
		mh.data = []int{}
		return out, nil
	}
	mh.data[0] = mh.data[mh.length]
	mh.heapifyDown(0)
	return out, nil
}

func (mh *MinHeap) parent(idx int) int {
	return (idx - 1) / 2
}

func (mh *MinHeap) leftChild(idx int) int {
	return 2*idx + 1
}

func (mh *MinHeap) rightChild(idx int) int {
	return 2*idx + 2
}

func (mh *MinHeap) heapifyUp(idx int) {
	if idx == 0 {
		return
	}

	pIdx := mh.parent(idx)
	pVal := mh.data[pIdx]
	val := mh.data[idx]

	if pVal > val {
		mh.data[idx] = pVal
		mh.data[pIdx] = val
		mh.heapifyUp(pIdx)
	}
}

func (mh *MinHeap) heapifyDown(idx int) {

	lIdx := mh.leftChild(idx)
	rIdx := mh.rightChild(idx)

	if lIdx >= mh.length || rIdx >= mh.length {
		return
	}

	lVal := mh.data[lIdx]
	rVal := mh.data[rIdx]
	val := mh.data[idx]

	if lVal > rVal && val > rVal {
		mh.data[idx] = rVal
		mh.data[rIdx] = val
		mh.heapifyDown(rIdx)
	} else if rVal > lVal && val > lVal {
		mh.data[idx] = lVal
		mh.data[lIdx] = val
		mh.heapifyDown(lIdx)
	}
}
