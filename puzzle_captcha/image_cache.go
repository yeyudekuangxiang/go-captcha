package captcha

import (
	"bytes"
	"image"
	"io/fs"
	"io/ioutil"
	"path"
	"strings"
)

var (
	bgImgCache [][]byte //缓存背景图片
	bkImgCache [][]byte //缓存滑块模板图片
)

type Images interface {
	Images() (files [][]byte, err error)
}
type FS struct {
	fsys fs.FS
	name string
}

func NewFS(fsys fs.FS, name string) *FS {
	return &FS{fsys: fsys, name: name}
}
func (f *FS) Images() (files [][]byte, err error) {
	dirs, err := fs.ReadDir(f.fsys, f.name)
	if err != nil {
		return nil, err
	}

	var fileArr [][]byte

	for _, d := range dirs {
		if d.IsDir() {
			continue
		}
		if strings.HasSuffix(d.Name(), ".png") {
			buf, err := fs.ReadFile(f.fsys, path.Join(f.name, d.Name()))
			if err != nil {
				return nil, err
			}
			fileArr = append(fileArr, buf)
		}
	}
	return fileArr, nil
}

type Path struct {
	p string
}

func NewPath(p string) *Path {
	return &Path{p: p}
}

func (f *Path) Images() (files [][]byte, err error) {
	return loadImages(f.p)
}

func LoadBackgroudImages(images Images) (err error) {
	bgImgCache, err = images.Images()
	return
}

func LoadBlockImages(images Images) (err error) {
	bkImgCache, err = images.Images()
	return
}

func loadImages(basePath string) ([][]byte, error) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}
	var fileArr [][]byte
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".png") {
			buf, err := ioutil.ReadFile(path.Join(basePath, f.Name()))
			if err != nil {
				return nil, err
			}
			fileArr = append(fileArr, buf)
		}
	}
	return fileArr, nil
}

// randBackgroudImage 随机抽取 背景图
func randBackgroudImage() (*ImageBuf, error) {
	n := r.Intn(len(bgImgCache))
	buf := bgImgCache[n]
	im, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	return &ImageBuf{
		w: im.Bounds().Dx(),
		h: im.Bounds().Dy(),
		i: im,
	}, nil
}

// randBlockImage 随机抽取 滑块图，和干扰图
func randBlockImage() (a *ImageBuf, b *ImageBuf, err error) {
	l := len(bkImgCache)
	n := r.Intn(len(bkImgCache))
	buf := bkImgCache[n]
	im, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, nil, err
	}
	var next = n + 1
	if next == l {
		next = 0
	}
	buf2 := bkImgCache[next]
	im2, _, err := image.Decode(bytes.NewReader(buf2))
	if err != nil {
		return nil, nil, err
	}
	a = &ImageBuf{
		w: im.Bounds().Dx(),
		h: im.Bounds().Dy(),
		i: im,
	}
	b = &ImageBuf{
		w: im2.Bounds().Dx(),
		h: im2.Bounds().Dy(),
		i: im2,
	}
	return
}
