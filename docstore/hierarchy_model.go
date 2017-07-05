package docstore

type Label string

type Labels []Label

type HierarchyLink struct {
	ParentLink  *HierarchyLink
	DocumentMap DocumentMap
	Label       Label
	ChildLink   *HierarchyLink
}

type HierarchyLinkyMap map[Label]*HierarchyLink
