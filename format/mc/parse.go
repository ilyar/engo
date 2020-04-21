package mc

import (
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

func parseMC(url string, r io.Reader) (*MovieClipResource, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	dataEgret, err := Unmarshal(content)
	if err != nil {
		return nil, err
	}

	textureName := path.Join(path.Dir(url), dataEgret.File)
	err = engo.Files.Load(textureName)
	if err != nil {
		return nil, err
	}

	spriteSheet, err := newSpritesheetFromFile(textureName, dataEgret)
	if err != nil {
		return nil, err
	}

	resource := &MovieClipResource{
		url:         url,
		SpriteSheet: spriteSheet,
		Actions:     make([]*common.Animation, 0),
	}

	for _, item := range dataEgret.Mc {
		for _, label := range item.Labels {
			action := &common.Animation{
				Name:   label.Name,
				Frames: make([]*common.Frame, 0, label.FrameEnd-label.FrameStart+1),
			}

			index := label.FrameStart - 1
			for _, frame := range item.Frames[label.FrameStart-1 : label.FrameEnd-1] {
				action.Frames = append(action.Frames, &common.Frame{
					Index: index,
					Bias: &engo.Point{
						X: float32(frame.X),
						Y: float32(frame.Y),
					},
				})
				index++
			}

			resource.Actions = append(resource.Actions, action)

			if strings.Contains(label.Name, "idle") {
				resource.DefaultAction = action
			}
		}
	}

	return resource, nil
}

func newSpritesheetFromFile(textureName string, mc MovieClip) (*common.Spritesheet, error) {
	res, err := engo.Files.Resource(textureName)
	if err != nil {
		return nil, fmt.Errorf("[MovieClip] [newSpritesheet] Resource for %q: %s", textureName, err.Error())
	}

	img, ok := res.(common.TextureResource)
	if !ok {
		return nil, fmt.Errorf("[MovieClip] [newSpritesheet] Resource not of type `TextureResource` for %q", textureName)
	}

	spriteRegions := make([]common.SpriteRegion, 0, len(mc.Regions))

	for _, frame := range mc.AllFrames() {
		region, exist := mc.Regions[frame.ResourceName]
		if !exist {
			continue
		}

		spriteRegions = append(spriteRegions, common.SpriteRegion{
			Position: engo.Point{float32(region.X), float32(region.Y)},
			Width:    region.W,
			Height:   region.H,
		})
	}

	return common.NewAsymmetricSpritesheetFromTexture(&img, spriteRegions), nil
}
