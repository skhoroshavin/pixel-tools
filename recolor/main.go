package main

import (
	"fmt"
	"path/filepath"
	"recolor/lut"
	"recolor/util"
)

func main() {
	cfg := util.LoadConfig("recolor-config.yaml")

	// Build LUT
	l := lut.New()
	ref := cfg.Reference
	if ref.BaseLUT != "" {
		fmt.Println("Using base LUT", ref.BaseLUT)
		l.Load(cfg.Reference.BaseLUT)
	}
	if ref.Folder != "" {
		for _, srcFile := range util.FindPairs(ref.Folder, ref.OriginalSuffix, ref.RecoloredSuffix) {
			dstFile := util.ReplaceSuffix(srcFile, ref.OriginalSuffix, ref.RecoloredSuffix)
			fmt.Println("Using reference images", filepath.Base(srcFile), "and", filepath.Base(dstFile))
			l.AddImageMapping(util.ReadImage(srcFile), util.ReadImage(dstFile))
		}
	}
	if ref.ResultingLUT != "" {
		fmt.Println("Saving resulting LUT as", ref.ResultingLUT)
		l.Save(ref.ResultingLUT)
	}

	// Recolor images
	rec := cfg.Recolor
	if rec.Folder != "" {
		for _, srcFile := range util.FindFiles(rec.Folder, rec.OriginalSuffix) {
			dstFile := util.ReplaceSuffix(srcFile, rec.OriginalSuffix, rec.RecoloredSuffix)
			fmt.Println("Applying LUT to image", srcFile, "saving as", dstFile)
			srcImage := util.ReadImage(srcFile)
			dstImage := l.ApplyToImage(srcImage)
			util.WriteImage(dstImage, dstFile)
		}
	}
}
