package data

type TriMesh struct {
	Vao                     Vao
	Vbo, Ibo                Buffer
	IndexCount, VertexCount int
}
