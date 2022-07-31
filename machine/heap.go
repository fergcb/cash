package machine

import (
	"cash/word"
	"fmt"
	"sort"
)

type Heap struct {
	Size  int
	Data  []word.Word
	Slots []Slot
}

type Slot struct {
	start int
	end   int
}

func (slot *Slot) length() int {
	return (slot.end - slot.start) + 1
}

func NewHeap(size int) *Heap {
	return &Heap{
		Size:  size,
		Data:  make([]word.Word, size),
		Slots: []Slot{{0, size}},
	}
}

func (heap *Heap) Allocate(size int) (int, error) {
	for _, slot := range heap.Slots {
		if slot.length() < size {
			continue
		}
		offset := slot.start
		slot.start += size
		heap.collapseSlots()
		return offset, nil
	}

	return -1, fmt.Errorf("not enough free memory on heap")
}

func (heap *Heap) Deallocate(offset int, size int) {
	slot := Slot{offset, offset + size - 1}
	heap.Slots = append(heap.Slots, slot)
	heap.collapseSlots()
}

func (heap *Heap) Write(offset int, size int, data []word.Word) int {
	for i := 0; i < size; i++ {
		heap.Data[offset+i] = data[i]
	}

	return offset + size
}

func (heap *Heap) Read(offset int, size int) []word.Word {
	data := make([]word.Word, size)
	for i := 0; i < size; i++ {
		data[i] = heap.Data[offset+i]
	}
	return data
}

func (heap *Heap) collapseSlots() {
	if len(heap.Slots) == 0 {
		return
	}

	sort.SliceStable(heap.Slots, func(i int, j int) bool {
		return heap.Slots[i].start < heap.Slots[j].start
	})

	stack := []Slot{heap.Slots[0]}
	top := 0

	for _, curSlot := range heap.Slots {
		topSlot := stack[top]
		if topSlot.end < curSlot.start {
			stack = append(stack, curSlot)
		} else if topSlot.end < curSlot.end {
			topSlot.end = curSlot.end
		}
	}

	heap.Slots = stack
}
