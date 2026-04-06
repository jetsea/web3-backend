// Package leetcode contains solutions to LeetCode problems implemented in Go.
package leetcode

// TwoSum returns the indices of the two numbers that add up to target.
// LeetCode #1 — https://leetcode.com/problems/two-sum/
//
// Approach: one-pass hash map.
//   - For each number, check whether (target - number) was already seen.
//   - If yes, return the stored index and the current index.
//   - If no, store the current number → index mapping.
//
// Time:  O(n)
// Space: O(n)
func TwoSum(nums []int, target int) []int {
	seen := make(map[int]int, len(nums)) // value → index

	for i, num := range nums {
		complement := target - num
		if j, ok := seen[complement]; ok {
			return []int{j, i}
		}
		seen[num] = i
	}

	return nil // no solution found (problem guarantees exactly one solution)
}
