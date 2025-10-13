package persist

import (
	"errors"
	"os"
	"path/filepath"

	"sewworkspacemanager/internal/domain"

	toml "github.com/pelletier/go-toml/v2"
)

type tomlNode struct {
	Label       string     `toml:"label"`
	Description string     `toml:"description,omitempty"`
	Command     string     `toml:"command,omitempty"`
	Args        []string   `toml:"args,omitempty"`
	WorkDir     string     `toml:"workdir,omitempty"`
	Children    []tomlNode `toml:"children,omitempty"`
}

type tomlMultiRoot struct {
	Nodes []tomlNode `toml:"nodes"`
}

type tomlDualGroup struct {
	Befehle  []tomlNode `toml:"befehle,omitempty"`
	Ablaeufe []tomlNode `toml:"ablaeufe,omitempty"`
}

type tomlTripleGroup struct {
	Befehle  []tomlNode `toml:"befehle,omitempty"`
	Ablaeufe []tomlNode `toml:"ablaeufe,omitempty"`
	Status   []tomlNode `toml:"status,omitempty"`
}

// toCommandNode konvertiert einen TOML-Knoten rekursiv in einen CommandNode.
func (n tomlNode) toCommandNode(parent *domain.CommandNode) *domain.CommandNode {
	cn := &domain.CommandNode{Label: n.Label, Description: n.Description, Command: n.Command, Args: n.Args, WorkDir: n.WorkDir, Parent: parent}
	for i := range n.Children {
		child := n.Children[i].toCommandNode(cn)
		cn.Children = append(cn.Children, child)
	}
	return cn
}

// fromCommandNode erzeugt die TOML-Repräsentation eines CommandNode rekursiv.
func fromCommandNode(n *domain.CommandNode) tomlNode {
	tn := tomlNode{Label: n.Label, Description: n.Description, Command: n.Command, Args: n.Args, WorkDir: n.WorkDir}
	if len(n.Children) > 0 {
		tn.Children = make([]tomlNode, len(n.Children))
		for i, ch := range n.Children {
			tn.Children[i] = fromCommandNode(ch)
		}
	}
	return tn
}

// LoadCommandTreeFromTOML lädt einen Befehlsbaum aus verschiedenen möglichen TOML-Layouts.
// Unterstützt folgende Varianten:
// - Triple Group (befehle + ablaeufe + status)
// - Dual Group (befehle + ablaeufe)
// - Multi Root (nodes[])
// - Einzelner Root-Knoten
// Alle Varianten werden in einen synthetischen Root (Label "") eingespannt, dessen Kinder die eigentlichen Wurzeln darstellen.
func LoadCommandTreeFromTOML(path string) (*domain.CommandNode, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var triple tomlTripleGroup
	if err := toml.Unmarshal(data, &triple); err == nil && (len(triple.Befehle) > 0 || len(triple.Ablaeufe) > 0 || len(triple.Status) > 0) {
		synth := &domain.CommandNode{Label: "", Expanded: true}
		for i := range triple.Befehle {
			synth.Children = append(synth.Children, triple.Befehle[i].toCommandNode(synth))
		}
		for i := range triple.Ablaeufe {
			synth.Children = append(synth.Children, triple.Ablaeufe[i].toCommandNode(synth))
		}
		for i := range triple.Status {
			synth.Children = append(synth.Children, triple.Status[i].toCommandNode(synth))
		}
		for _, ch := range synth.Children {
			if depthLimitExceeded(ch, 0) {
				return nil, errors.New("Befehlsbaum überschreitet maximale Tiefe 3")
			}
		}
		flattenOneLevel(synth)
		linkParents(synth)
		return synth, nil
	}
	var dual tomlDualGroup
	if err := toml.Unmarshal(data, &dual); err == nil && (len(dual.Befehle) > 0 || len(dual.Ablaeufe) > 0) {
		synth := &domain.CommandNode{Label: "", Expanded: true}
		for i := range dual.Befehle {
			synth.Children = append(synth.Children, dual.Befehle[i].toCommandNode(synth))
		}
		for i := range dual.Ablaeufe {
			synth.Children = append(synth.Children, dual.Ablaeufe[i].toCommandNode(synth))
		}
		for _, ch := range synth.Children {
			if depthLimitExceeded(ch, 0) {
				return nil, errors.New("Befehlsbaum überschreitet maximale Tiefe 3")
			}
		}
		flattenOneLevel(synth)
		linkParents(synth)
		return synth, nil
	}
	var multi tomlMultiRoot
	if err := toml.Unmarshal(data, &multi); err == nil && len(multi.Nodes) > 0 {
		synth := &domain.CommandNode{Label: "", Expanded: true}
		for i := range multi.Nodes {
			synth.Children = append(synth.Children, multi.Nodes[i].toCommandNode(synth))
		}
		for _, ch := range synth.Children {
			if depthLimitExceeded(ch, 0) {
				return nil, errors.New("Befehlsbaum überschreitet maximale Tiefe 3")
			}
		}
		flattenOneLevel(synth)
		linkParents(synth)
		return synth, nil
	}
	var rootTN tomlNode
	if err := toml.Unmarshal(data, &rootTN); err != nil {
		return nil, err
	}
	root := rootTN.toCommandNode(nil)
	if depthLimitExceeded(root, 0) {
		return nil, errors.New("Befehlsbaum überschreitet maximale Tiefe 3")
	}
	flattenOneLevel(root)
	linkParents(root)
	return root, nil
}

// SaveCommandTreeToTOML speichert den Baum. Bei leerem Root-Label wird ein Multi-Root Layout erzeugt.
func SaveCommandTreeToTOML(root *domain.CommandNode, path string) error {
	if root == nil {
		return errors.New("root is nil")
	}
	var data []byte
	var err error
	if root.Label == "" {
		multi := tomlMultiRoot{}
		for _, ch := range root.Children {
			multi.Nodes = append(multi.Nodes, fromCommandNode(ch))
		}
		data, err = toml.Marshal(multi)
	} else {
		data, err = toml.Marshal(fromCommandNode(root))
	}
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// SaveDualCommandTrees speichert zwei getrennte Gruppen (befehle + ablaeufe).
func SaveDualCommandTrees(befehleRoot, ablaeufeRoot *domain.CommandNode, path string) error {
	dual := tomlDualGroup{}
	if befehleRoot != nil {
		for _, ch := range befehleRoot.Children {
			dual.Befehle = append(dual.Befehle, fromCommandNode(ch))
		}
	}
	if ablaeufeRoot != nil {
		for _, ch := range ablaeufeRoot.Children {
			dual.Ablaeufe = append(dual.Ablaeufe, fromCommandNode(ch))
		}
	}
	data, err := toml.Marshal(dual)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// SaveAllCommandGroups speichert alle drei Gruppen (befehle + ablaeufe + status).
func SaveAllCommandGroups(befehleRoot, ablaeufeRoot, statusRoot *domain.CommandNode, path string) error {
	triple := tomlTripleGroup{}
	if befehleRoot != nil {
		for _, ch := range befehleRoot.Children {
			triple.Befehle = append(triple.Befehle, fromCommandNode(ch))
		}
	}
	if ablaeufeRoot != nil {
		for _, ch := range ablaeufeRoot.Children {
			triple.Ablaeufe = append(triple.Ablaeufe, fromCommandNode(ch))
		}
	}
	if statusRoot != nil {
		for _, ch := range statusRoot.Children {
			triple.Status = append(triple.Status, fromCommandNode(ch))
		}
	}
	data, err := toml.Marshal(triple)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadAllCommandGroups lädt alle drei Gruppen, versucht zuerst das Triple-Format,
// fällt zurück auf Dual-Format oder Legacy-Einzelbaum und erzeugt leere Platzhalter für fehlende Gruppen.
func LoadAllCommandGroups(path string) (befehleRoot, ablaeufeRoot, statusRoot *domain.CommandNode, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil, err
	}
	var triple tomlTripleGroup
	if err := toml.Unmarshal(data, &triple); err == nil && (len(triple.Befehle) > 0 || len(triple.Ablaeufe) > 0 || len(triple.Status) > 0) {
		br := &domain.CommandNode{Label: "", Expanded: true}
		for i := range triple.Befehle {
			br.Children = append(br.Children, triple.Befehle[i].toCommandNode(br))
		}
		ar := &domain.CommandNode{Label: "", Expanded: true}
		for i := range triple.Ablaeufe {
			ar.Children = append(ar.Children, triple.Ablaeufe[i].toCommandNode(ar))
		}
		sr := &domain.CommandNode{Label: "", Expanded: true}
		for i := range triple.Status {
			sr.Children = append(sr.Children, triple.Status[i].toCommandNode(sr))
		}
		flattenOneLevel(br)
		linkParents(br)
		flattenOneLevel(ar)
		linkParents(ar)
		flattenOneLevel(sr)
		linkParents(sr)
		return br, ar, sr, nil
	}
	br, ar, derr := LoadDualCommandTrees(path)
	if derr == nil {
		return br, ar, &domain.CommandNode{Label: "", Expanded: true}, nil
	}
	legacy, lerr := LoadCommandTreeFromTOML(path)
	if lerr != nil {
		return nil, nil, nil, lerr
	}
	return legacy, &domain.CommandNode{Label: "", Expanded: true}, &domain.CommandNode{Label: "", Expanded: true}, nil
}

// LoadDualCommandTrees lädt nur die zwei Gruppen (befehle + ablaeufe) oder fällt auf Legacy zurück.
func LoadDualCommandTrees(path string) (befehleRoot, ablaeufeRoot *domain.CommandNode, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	var dual tomlDualGroup
	if err := toml.Unmarshal(data, &dual); err == nil && (len(dual.Befehle) > 0 || len(dual.Ablaeufe) > 0) {
		br := &domain.CommandNode{Label: "", Expanded: true}
		for i := range dual.Befehle {
			br.Children = append(br.Children, dual.Befehle[i].toCommandNode(br))
		}
		ar := &domain.CommandNode{Label: "", Expanded: true}
		for i := range dual.Ablaeufe {
			ar.Children = append(ar.Children, dual.Ablaeufe[i].toCommandNode(ar))
		}
		flattenOneLevel(br)
		linkParents(br)
		flattenOneLevel(ar)
		linkParents(ar)
		return br, ar, nil
	}
	root, lerr := LoadCommandTreeFromTOML(path)
	if lerr != nil {
		return nil, nil, lerr
	}
	return root, &domain.CommandNode{Label: "", Expanded: true}, nil
}

// depthLimitExceeded validiert die maximale erlaubte Tiefe (3 Ebenen inklusive Root-Kinder) rekursiv.
func depthLimitExceeded(n *domain.CommandNode, depth int) bool {
	if depth > 2 {
		return true
	}
	for _, ch := range n.Children {
		if depthLimitExceeded(ch, depth+1) {
			return true
		}
	}
	return false
}

// linkParents setzt Parent-Zeiger rekursiv.
func linkParents(n *domain.CommandNode) {
	for _, ch := range n.Children {
		ch.Parent = n
		linkParents(ch)
	}
}

// flattenOneLevel entfernt reine Container-Knoten ohne eigenen Befehl, sofern deren Kinder ebenfalls reine Container sind,
// um eine unnötige Zwischenebene zu vermeiden.
func flattenOneLevel(root *domain.CommandNode) {
	if root == nil {
		return
	}
	changed := true
	for changed {
		changed = false
		var newChildren []*domain.CommandNode
		for _, ch := range root.Children {
			pureContainer := ch.Command == "" && len(ch.Children) > 0
			if pureContainer {
				hasCommandChild := false
				for _, g := range ch.Children {
					if g.Command != "" {
						hasCommandChild = true
						break
					}
				}
				if !hasCommandChild {
					for _, g := range ch.Children {
						g.Parent = root
						newChildren = append(newChildren, g)
					}
					changed = true
					continue
				}
			}
			newChildren = append(newChildren, ch)
		}
		if changed {
			root.Children = newChildren
		}
	}
}
