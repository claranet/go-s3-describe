package main

import (
	"fmt"
	"os"
	"strconv"

	"code.cloudfoundry.org/bytefmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

func PrintResult(s3Buckets *[]s3Bucket) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Bucket", "Region", "IsPublic", "NumberOfObjects", "Standard (BYTES)", "StandardIA (BYTES)", "ReducedRedundancy (BYTES)", "Glacier (BYTES)"})
	for _, v := range *s3Buckets {
		// Convert BucketSize from float64 to a human readable format
		t.AppendRow(table.Row{
			v.name,
			v.region,
			strconv.FormatBool(v.isPublic),
			fmt.Sprintf("%v", v.numberOfObjects),
			bytefmt.ByteSize(uint64(v.bucketSizeBytes["StandardStorage"])),
			bytefmt.ByteSize(uint64(v.bucketSizeBytes["StandardIAStorage"])),
			bytefmt.ByteSize(uint64(v.bucketSizeBytes["ReducedRedundancyStorage"])),
			bytefmt.ByteSize(uint64(v.bucketSizeBytes["GlacierStorage"]))})
	}
	t.SetRowPainter(table.RowPainter(func(row table.Row) text.Colors {
		if row[2] == "true" {
			return text.Colors{text.BgRed, text.FgWhite}
		}
		if row[3] == "0" {
			return text.Colors{text.BgYellow, text.FgBlack}
		}
		return nil
	}))
	t.SetStyle(table.StyleColoredBright)
	t.Render()
	return nil
}
