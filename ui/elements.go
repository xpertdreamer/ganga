package ui

import (
	"fmt"
	"ganga/back"
	"strings"
)

func CreatePreview(k* back.KBase) string {
	var sb strings.Builder
	sb.Grow(512)

	sb.WriteString("Objects : \n\n")
	for _, obj := range k.Objects {
		sb.WriteString("\t")
		sb.WriteString(obj.Description)
		sb.WriteString(" ")
		fmt.Fprintf(&sb, "%v", obj.Properties)
		sb.WriteString(",\n\n")
	}

	sb.WriteString("Properties : \n\n")
	for _, prop := range k.Properties {
		sb.WriteString("\t")
		sb.WriteString(prop.ID)
		sb.WriteString(":")
		sb.WriteString(prop.Description)
		sb.WriteString(",\n")
	}

	return sb.String()
}
