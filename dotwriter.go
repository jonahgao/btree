package btree

import (
	"bytes"
	"fmt"
	"os"
)

func writeDot(root *node, file string) error {
	buffer := bytes.NewBuffer(nil)

	//
	_, err := buffer.WriteString(
		`
digraph {
    graph [margin=0, splines=line];
    edge [penwidth=2];
    node [shape = record, penwidth=2, style=filled, fillcolor=white];

`)
	if err != nil {
		return nil
	}

	startIndex := 0
	err = writeNode(buffer, root, &startIndex)

	_, err = buffer.WriteString("\n}")
	if err != nil {
		return err
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	_, err = f.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return nil
	}

	return f.Close()
}

func writeNode(buf *bytes.Buffer, node *node, startIdx *int) error {
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
		err = writeNode(buf, c, startIdx)
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
