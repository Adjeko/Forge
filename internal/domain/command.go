package domain

import "strconv"

// CommandNode repräsentiert entweder einen Ordner (Kinder) oder einen konkreten Befehl (ausführbare Details).
type CommandNode struct {
	Label       string
	Description string
	Command     string
	Args        []string
	WorkDir     string
	Children    []*CommandNode
	Parent      *CommandNode
	Expanded    bool
}

// DepthOf liefert die logische Tiefe ohne den synthetischen Root.
func DepthOf(n *CommandNode) int {
	if n == nil {
		return 0
	}
	d := 0
	for n.Parent != nil {
		if n.Parent.Parent == nil { // Parent ist synthetischer Root
			return d
		}
		d++
		n = n.Parent
	}
	return d
}

// PadRight füllt s mit Leerzeichen bis Breite w.
func PadRight(s string, w int) string {
	r := []rune(s)
	if len(r) >= w {
		return s
	}
	return s + string(make([]rune, w-len(r)))
}

// Itoa konvertiert int zu string.
func Itoa(i int) string { return fmtInt(i) }

// fmtInt kapselt strconv.Itoa ohne es überall zu importieren.
func fmtInt(i int) string { return strconv.Itoa(i) }

// AddNodeToTree fügt node unter dem gegebenen folderPath ein.
func AddNodeToTree(root *CommandNode, folderPath []string, node *CommandNode) {
	if root == nil || node == nil {
		return
	}
	idx := 0
	if len(folderPath) > 0 && (folderPath[0] == "(Top)" || folderPath[0] == root.Label) {
		idx = 1
	}
	cur := root
	depth := 0
	for ; idx < len(folderPath); idx++ {
		seg := folderPath[idx]
		var found *CommandNode
		for _, ch := range cur.Children {
			if ch.Label == seg {
				found = ch
				break
			}
		}
		if found == nil {
			found = &CommandNode{Label: seg, Expanded: true, Parent: cur}
			cur.Children = append(cur.Children, found)
		}
		cur = found
		depth++
		if depth >= 1 { // nur eine Ordner-Ebene erlauben (Variante A)
			break
		}
	}
	node.Parent = cur
	cur.Children = append(cur.Children, node)
}

// (entfernte doppelte englische Blockdefinitionen)
