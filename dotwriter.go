package btree

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func writeDotSvg(dotExePath string, outputSvg string, tree *Btree) error {
	buffer := bytes.NewBuffer(nil)
	err := writeDotGraph(tree.root, buffer)
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

func writeDotGraph(root *node, buffer *bytes.Buffer) error {
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
		err = writeDotNode(buffer, root, &startIndex)
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

func writeDotNode(buf *bytes.Buffer, node *node, startIdx *int) error {
	prevIdx := *startIdx
	nodeStr := fmt.Sprintf("    node%d[label= \"<f0> ● ", *startIdx)
	for i := 0; i < node.numKeys; i++ {
		nodeStr = nodeStr + fmt.Sprintf("| %v | <f%d> ● ", string(node.keys[i]), i+1)
	}
	nodeStr = nodeStr + "\"]\n"
	*startIdx = *startIdx + 1

	_, err := buf.WriteString(nodeStr)
	if err != nil {
		return err
	}

	for i, c := range node.children {
		err = writeDotNode(buf, c, startIdx)
		if err != nil {
			return err
		}

		linkStr := fmt.Sprintf("    node%d:f%d -> node%d\n", prevIdx, i, *startIdx-1)
		_, err = buf.WriteString(linkStr)
		if err != nil {
			return err
		}
	}
	return nil
}
