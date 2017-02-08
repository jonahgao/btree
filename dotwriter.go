package btree

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// dump btree to svg picture (use graphviz)
func writeDotSvg(dotExePath string, outputSvg string, tree *MVCCBtree) error {
	buffer := bytes.NewBuffer(nil)
	err := writeDotGraph(tree.GetTree().root, buffer)
	if err != nil {
		return err
	}

	cmd := exec.Command(dotExePath, "-Tsvg")
	cmd.Stdin = buffer
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(outputSvg, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	_, err = f.Write(out.Bytes())
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return nil
	}

	return f.Close()
}

func writeDotGraph(root node, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(
		`
digraph {
    graph [margin=0, splines=line];
    edge [penwidth=2];
    node [shape = record,style=filled, fillcolor=white];

`)
	if err != nil {
		return nil
	}

	if root != nil {
		startIndex := 0
		_, err = writeDotNode(buffer, root, &startIndex)
		if err != nil {
			return err
		}
	}

	_, err = buffer.WriteString("\n}")
	if err != nil {
		return err
	}
	return nil
}

func writeDotNode(buf *bytes.Buffer, node node, startIdx *int) (nodeIndex int, err error) {
	nodeIndex = *startIdx
	nodeStr := fmt.Sprintf("    node%d[label= \"<f0> ● ", nodeIndex)
	for i := 0; i < node.numOfKeys(); i++ {
		nodeStr = nodeStr + fmt.Sprintf("| %v | <f%d> ● ", string(node.keyAt(i)), i+1)
	}
	nodeStr = nodeStr + "\"]\n"
	*startIdx = *startIdx + 1

	_, err = buf.WriteString(nodeStr)
	if err != nil {
		return
	}

	if !node.isLeaf() {
		n := node.(*internalNode)
		for i, c := range n.children {
			var idx int
			idx, err = writeDotNode(buf, c, startIdx)
			if err != nil {
				return
			}

			linkStr := fmt.Sprintf("    node%d:f%d -> node%d\n", nodeIndex, i, idx)
			_, err = buf.WriteString(linkStr)
			if err != nil {
				return
			}
		}
	}
	return
}
