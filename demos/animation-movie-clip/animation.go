//+build demo

package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/engo/format/mc"
)

const (
	Height = 700
	Width  = Height * 1.6
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	opts := engo.RunOptions{
		Title: "Animation Demo",

		FPSLimit: 25,
		Width:    Width,
		Height:   Height,

		//HeadlessMode: true,
		StandardInputs: true,
	}
	engo.Run(opts, &DefaultScene{})
}

func NewHeroEntity(position, scale engo.Point, mcr *mc.MovieClipResource) *HeroEntity {
	entity := &HeroEntity{BasicEntity: ecs.NewBasic()}

	entity.RenderComponent = common.RenderComponent{}
	entity.RenderComponent.Drawable = mcr.Drawable
	entity.RenderComponent.Scale = scale

	entity.SpaceComponent = common.SpaceComponent{}
	entity.SpaceComponent.Width = mcr.Drawable.Width()
	entity.SpaceComponent.Height = mcr.Drawable.Height()
	entity.SpaceComponent.Position = engo.Point{
		X: position.X + mcr.Drawable.Width()*scale.X,
		Y: position.Y - mcr.Drawable.Height()*scale.Y,
	}

	entity.AnimationComponent = common.NewAnimationComponent(mcr.SpriteSheet.Drawables(), 0.0)
	entity.AnimationComponent.AddAnimations(mcr.Actions)
	entity.AnimationComponent.AddDefaultAnimation(mcr.DefaultAction)

	return entity
}

type HeroEntity struct {
	ecs.BasicEntity
	common.AnimationComponent
	common.RenderComponent
	common.SpaceComponent
}

type DefaultScene struct{}

func (*DefaultScene) Preload() {
	err := engo.Files.Load(
		"sheep.mc.json",
		"Engo.mc.json",
	)
	if err != nil {
		log.Fatalln(err)
	}
}

func (scene *DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.Alpha16{A: 0x7575})

	w.AddSystemInterface(&common.RenderSystem{}, new(common.Renderable), nil)
	w.AddSystemInterface(&common.AnimationSystem{}, new(common.Animationable), nil)
	w.AddSystemInterface(&ControlSystem{}, new(Controllable), nil)

	baseLine := float32(Height - 100)
	w.AddEntity(NewFieldEntity(
		engo.Point{X: 0, Y: baseLine},
		engo.Point{X: Width, Y: 100},
		color.Black,
	))

	position := engo.Point{X: 200, Y: baseLine + 20}
	mcr, err := mc.LoadResource("Engo.mc.json")
	if err != nil {
		log.Fatalln(err)
	}
	hero := NewHeroEntity(position, engo.Point{X: 1, Y: 1}, mcr)
	w.AddEntity(hero)
}

func (*DefaultScene) Type() string { return "GameWorld" }

type controlEntity struct {
	*ecs.BasicEntity
	*common.AnimationComponent
}

type Controllable interface {
	common.BasicFace
	common.AnimationFace
}

type ControlSystem struct {
	entities map[uint64]controlEntity
}

func (c *ControlSystem) AddByInterface(i ecs.Identifier) {
	o, _ := i.(Controllable)
	c.Add(o.GetBasicEntity(), o.GetAnimationComponent())
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, anim *common.AnimationComponent) {
	if c.entities == nil {
		c.entities = make(map[uint64]controlEntity)
	}
	c.entities[basic.ID()] = controlEntity{basic, anim}
}

func (c *ControlSystem) Remove(basic ecs.BasicEntity) {
	if c.entities != nil {
		delete(c.entities, basic.ID())
	}
}

func (c *ControlSystem) Update(dt float32) {
	for _, e := range c.entities {
		if engo.Input.Button("action").JustPressed() {
			c.randAction(e.GetAnimationComponent())
		}
	}
}

func (c *ControlSystem) randAction(anim *common.AnimationComponent) {
	animCount := len(anim.Animations)
	list := make([]string, 0, animCount)
	for name := range anim.Animations {
		list = append(list, name)
	}
	anim.SelectAnimationByName(list[rand.Intn(animCount)])
}

type FieldEntity struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func NewFieldEntity(position, size engo.Point, color color.Color) *FieldEntity {
	entity := &FieldEntity{BasicEntity: ecs.NewBasic()}
	entity.SpaceComponent = common.SpaceComponent{Position: position}
	entity.SpaceComponent.Width = size.X
	entity.SpaceComponent.Height = size.Y
	entity.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}}
	entity.RenderComponent.Color = color

	return entity
}
