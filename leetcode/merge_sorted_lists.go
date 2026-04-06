package leetcode

// MergeTwoSortedLists merges two sorted linked lists and returns the merged list.
// LeetCode #21 — https://leetcode.com/problems/merge-two-sorted-lists/
//
// Approach: dummy-head pointer.
//   - Maintain a dummy head node so we never have to special-case the first element.
//   - At each step, pick the smaller front node and advance that list.
//   - Append the remaining non-nil list when one is exhausted.
//
// Time:  O(m + n)
// Space: O(1)  — no extra nodes created, only pointers rearranged
// ListNode is the singly-linked list node definition used by LeetCode.
type ListNode struct {
	Val  int
	Next *ListNode
}

// NewList builds a linked list from a slice of values.
func NewList(values []int) *ListNode {
	if len(values) == 0 {
		return nil
	}
	head := &ListNode{Val: values[0]}
	current := head
	for _, v := range values[1:] {
		current.Next = &ListNode{Val: v}
		current = current.Next
	}
	return head
}

// ToSlice converts a linked list to a slice.
func ToSlice(head *ListNode) []int {
	var result []int
	for head != nil {
		result = append(result, head.Val)
		head = head.Next
	}
	return result
}

func MergeTwoSortedLists(list1 *ListNode, list2 *ListNode) *ListNode {
	dummy := &ListNode{} // sentinel head — never returned
	current := dummy

	for list1 != nil && list2 != nil {
		if list1.Val <= list2.Val {
			current.Next = list1
			list1 = list1.Next
		} else {
			current.Next = list2
			list2 = list2.Next
		}
		current = current.Next
	}

	// Attach the remaining non-nil list.
	if list1 != nil {
		current.Next = list1
	} else {
		current.Next = list2
	}

	return dummy.Next
}
