package index

// finds all docs containing every term in the user query
// runs in O(n) where n is the len of shortest posting list

func intersect(a, b []uint32) []uint32 {
	var out []uint32
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

func (idx *Index) Search(query string) []Document {
	terms := Analyze(query)
	if len(terms) == 0 {
		return nil
	}

	lists := make([][]uint32, 0, len(terms))

	for _, term := range terms {
		p := idx.Postings(term)
		if len(p) == 0 {
			return nil
		}
		lists = append(lists, p)
	}
	sortByLen(lists)

	result := lists[0]
	for _, list := range lists[1:] {
		result = intersect(result, list)
		if len(result) == 0 {
			return nil
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

func sortByLen(lists [][]uint32) {

	for i := 1; i < len(lists); i++ {
		for j := i; j > 0 && len(lists[j]) < len(lists[j-1]); j-- {
			lists[j], lists[j-1] = lists[j-1], lists[j]
		}
	}
}
