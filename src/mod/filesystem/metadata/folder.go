package metadata

import (
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

/*
	Generate folder thumbnail from the containing files
	The preview is generated by overlapping 2 - 3 layers of images
*/

func generateThumbnailForFolder(cacheFolder string, file string, generateOnly bool) (string, error) {
	//Check if this folder has cache image folder
	cacheFolderInsideThisFolder := filepath.Join(file, "/.cache")
	if !fileExists(cacheFolderInsideThisFolder) {
		//This folder do not have a cache folder
		return "", errors.New("No previewable files")
	}

	//Load the base template
	if !fileExists("web/img/system/folder-preview.png") {
		//Missing system files. Skip rendering
		return "", errors.New("Missing system template image file")
	}
	image1, err := os.Open("web/img/system/folder-preview.png")
	if err != nil {
		return "", err
	}

	baseTemplate, err := png.Decode(image1)
	if err != nil {
		return "", err
	}
	image1.Close()

	//Generate the base image
	b := baseTemplate.Bounds()
	resultThumbnail := image.NewRGBA(b)
	draw.Draw(resultThumbnail, b, baseTemplate, image.ZP, draw.Over)

	//Get cached file inside this folder
	contentCache, _ := filepath.Glob(filepath.Join(cacheFolderInsideThisFolder, "/*.jpg"))

	//Check if there are more than 1 file inside this folder that is cached
	if len(contentCache) > 1 {
		//More than 1 files. Render the image at the back
		image2, err := os.Open(contentCache[1])
		if err != nil {
			return "", err
		}
		backImage, err := jpeg.Decode(image2)
		if err != nil {
			return "", err
		}
		backImgOffset := image.Pt(155, 110)
		defer image2.Close()
		resizedBackImg := resize.Resize(250, 250, backImage, resize.Lanczos3)
		draw.Draw(resultThumbnail, resizedBackImg.Bounds().Add(backImgOffset), resizedBackImg, image.ZP, draw.Over)
	} else if len(contentCache) == 0 {
		//Nothing to preview inside this folder
		return "", errors.New("No previewable files")
	}

	//Render the top image
	image3, err := os.Open(contentCache[0])
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}

	topImage, err := jpeg.Decode(image3)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	defer image3.Close()

	topImageOffset := image.Pt(210, 210)
	resizedTopImage := resize.Resize(260, 260, topImage, resize.Lanczos3)
	draw.Draw(resultThumbnail, resizedTopImage.Bounds().Add(topImageOffset), resizedTopImage, image.ZP, draw.Over)

	outfile, err := os.Create(filepath.Join(cacheFolder, filepath.Base(file)+".png"))
	if err != nil {
		log.Fatalf("failed to create: %s", err)
	}
	png.Encode(outfile, resultThumbnail)
	outfile.Close()

	ctx, err := getImageAsBase64(cacheFolder + filepath.Base(file) + ".png")
	return ctx, err
}
