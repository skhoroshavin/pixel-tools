package main

import (
	"fmt"
	"path/filepath"

	"pixel-tools/cmd/recolor/lut"
	"pixel-tools/cmd/recolor/util"
	"pixel-tools/pkg/imgutil"
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
		fmt.Println("Looking for reference image pairs in", ref.Folder)
		for _, srcFile := range util.FindPairs(ref.Folder, ref.OriginalSuffix, ref.RecoloredSuffix) {
			dstFile := util.ReplaceSuffix(srcFile, ref.OriginalSuffix, ref.RecoloredSuffix)
			fmt.Println("Using reference images", filepath.Base(srcFile), "and", filepath.Base(dstFile))
			l.AddImageMapping(imgutil.ReadImage(srcFile), imgutil.ReadImage(dstFile))
		}
	}
	if ref.ResultingLUT != "" {
		fmt.Println("Saving resulting LUT as", ref.ResultingLUT)
		l.Save(ref.ResultingLUT)
	}

	// Recolor images
	recolor := cfg.Recolor
	if recolor.Folder != "" {
		fmt.Println("Looking for images to recolor in", recolor.Folder)
		for _, srcFile := range util.FindFiles(recolor.Folder, recolor.OriginalSuffix) {
			dstFile := util.ReplaceSuffix(srcFile, recolor.OriginalSuffix, recolor.RecoloredSuffix)
			fmt.Println("Applying LUT to image", srcFile, "saving as", dstFile)
			srcImage := imgutil.ReadImage(srcFile)
			dstImage := l.ApplyToImage(srcImage)
			imgutil.WriteImage(dstImage, dstFile)
		}
	}
}
