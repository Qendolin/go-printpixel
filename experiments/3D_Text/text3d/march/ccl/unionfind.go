package ccl

type UnionFind struct {
	// Name of the next NextLabel, when one is created
	NextLabel int
	// Array which holds label -> set equivalences
	labels []int
}

func (uf *UnionFind) MakeLabel() int {
	l := uf.NextLabel
	uf.NextLabel++
	uf.labels = append(uf.labels, l)
	return l
}

// setRoot makes all nodes "in the path of node i" point to root
func (uf *UnionFind) setRoot(i, root int) {
	for uf.labels[i] < i {
		j := uf.labels[i]
		uf.labels[i] = root
		i = j
	}
	uf.labels[i] = root
}

// FindRoot finds the root node of the tree containing node i
func (uf *UnionFind) findRoot(i int) int {
	for uf.labels[i] < i {
		i = uf.labels[i]
	}
	return i
}

// Find finds the root of the tree containing node i
// Simultaneously compresses the tree
func (uf *UnionFind) Find(i int) int {
	root := uf.findRoot(i)
	uf.setRoot(i, root)
	return root
}

// Joins the two trees containing nodes i and j
// Modified to be less agressive about compressing paths
// because performance was suffering some from over-compression
func (uf *UnionFind) Union(i, j int) {
	if i == j {
		return
	}

	rooti := uf.findRoot(i)
	rootj := uf.findRoot(j)

	root := rooti
	if rooti > rootj {
		root = rootj
	}

	uf.setRoot(j, root)
	uf.setRoot(j, root)
}

func (uf *UnionFind) Flatten() {
	for i := 1; i < len(uf.labels); i++ {
		uf.labels[i] = uf.labels[uf.labels[i]]
	}
}
