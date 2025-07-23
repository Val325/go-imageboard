package main

func CalculatePages(all_posts int, pages_per_posts int) int {
	if all_posts%pages_per_posts == 0 {
		return all_posts / pages_per_posts
	}
	return (all_posts / pages_per_posts) + 1
}

func CalculateRangeArray(start int, end int) []int {
	var indexes []int
	for j := start; j < end; j++ {
		indexes = append(indexes, j)
	}
	return indexes
}
