package index

// Search returns all documents that contain every term in the query string.
// It uses a two-pointer posting-list intersection that runs in O(n) where
// n is the length of the shortest posting list.
func (idx *Index) Search(query string) []Document {
	terms := Analyze(query)
	if len(terms) == 0 {
		return []Document{}
	}

	lists := make([][]uint32, 0, len(terms))

	for _, term := range terms {
		p := idx.Postings(term)
		if len(p) == 0 {
			return []Document{}
		}
		lists = append(lists, p)
	}
	sortByLen(lists)

	result := lists[0]
	for _, list := range lists[1:] {
		result = intersect(result, list)
		if len(result) == 0 {
			return []Document{}
		}
	}

	docs := make([]Document, 0, len(result))
	for _, id := range result {
		if doc, ok := idx.GetDoc(id); ok {
			docs = append(docs, doc)
		}
	}

	return docs
}

// intersect returns the sorted IDs present in both a and b using a
// two-pointer merge. Both slices must be sorted in ascending order.
func intersect(a, b []uint32) []uint32 {
	out := make([]uint32, 0, len(a))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] == b[j]:
			out = append(out, a[i])
			i++
			j++

		case a[i] < b[j]:
			i++

		default:
			j++
		}
	}
	return out
}

func sortByLen(lists [][]uint32) {
	for i := 1; i < len(lists); i++ {
		for j := i; j > 0 && len(lists[j]) < len(lists[j-1]); j-- {
			lists[j], lists[j-1] = lists[j-1], lists[j]
		}
	}
}
