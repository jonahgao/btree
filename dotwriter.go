package btree

import (
	"os"
)

func writeDot(root *node, file string) error {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	_, err = f.WriteString(`
digraph {
    graph [margin=0, splines=line];
    edge [penwidth=2];
    node [shape = record, margin=0.03,1.2, penwidth=2, style=filled, fillcolor=white];
`)

	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return nil
	}

	return f.Close()
}
