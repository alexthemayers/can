package openapi

type Traversable interface {
	getChildren() map[string]Traversable
	setChild(i string, t Traversable)
	getParent() Traversable
	setParent(parent Traversable)
	GetName() string
	setName(name string)
	getBasePath() string
	getRef() string
	setRenderer(r Renderer)
	getRenderer() Renderer
}

type TraversalFunc func(key string, parent, child Traversable) (Traversable, error)

type node struct {
	basePath string
	parent   Traversable
	name     string
	renderer Renderer
	ref      string
}

var _ Traversable = &node{}

func (n *node) getChildren() map[string]Traversable {
	panic("not implemented by composed type")
}

func (n *node) setChild(i string, t Traversable) {
	panic("not implemented by composed type")
}

func (n *node) getParent() Traversable {
	return n.parent
}

func (n *node) setParent(parent Traversable) {
	n.parent = parent
}

func (n *node) getBasePath() string {
	return n.parent.getBasePath()
}

func (n *node) getRef() string {
	return n.ref
}

func (n *node) GetName() string {
	return n.renderer.sanitiseName(n.parent.GetName() + n.name)
}

func (n *node) setName(name string) {
	n.name = name
}

func (n *node) setRenderer(r Renderer) {
	n.renderer = r
}

func (n *node) getRenderer() Renderer {
	return n.renderer
}

// Traverse takes a Traversable node and applies some function to the node within the tree. It recursively calls itself and fails early when an error is thrown
func Traverse(node Traversable, f TraversalFunc) (Traversable, error) {
	if node == nil || f == nil {
		return node, nil
	}
	var recurse func(node Traversable, f TraversalFunc) (Traversable, error)
	recurse = func(node Traversable, f TraversalFunc) (Traversable, error) {
		children := node.getChildren()
		for i := range children {
			child := children[i]
			if child == nil {
				continue
			}
			// Update Child Node
			newChild, err := f(i, node, child)
			if err != nil {
				return nil, err
			}
			node.setChild(i, newChild)

			if newChild == nil {
				continue
			}
			_, err = recurse(newChild, f)
			if err != nil {
				return nil, err
			}
		}

		return node, nil
	}
	node, err := f("", nil, node)
	if err != nil {
		return nil, err
	}

	return recurse(node, f)
}

func Dig(node Traversable, key ...string) Traversable {
	if len(key) == 0 {
		return node
	}
	return Dig(node.getChildren()[key[0]], key[1:]...)
}
