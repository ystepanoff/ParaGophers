package game

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ystepanoff/paragopher/internal/config"
	"github.com/ystepanoff/paragopher/internal/utils"
)

type Helicopter struct {
	x, y        float32
	leftToRight bool
	lastDrop    time.Time
}

func (g *Game) drawHelicopter(screen *ebiten.Image, h *Helicopter) {
	vector.DrawFilledRect(
		screen,
		h.x-config.HelicopterBodyWidth/2.0,
		h.y-config.HelicopterBodyHeight/2.0,
		config.HelicopterBodyWidth,
		config.HelicopterBodyHeight,
		config.ColourTeal,
		false,
	)
	tailX := h.x - config.HelicopterBodyWidth
	if !h.leftToRight {
		tailX = h.x + config.HelicopterBodyWidth/2.0
	}
	vector.DrawFilledRect(
		screen,
		tailX,
		h.y-config.HelicopterTailHeight/2.0,
		config.HelicopterTailWidth,
		config.HelicopterTailHeight,
		config.ColourTeal,
		false,
	)
	vector.StrokeLine(
		screen,
		h.x-config.HelicopterRotorLen/2.0,
		h.y-config.HelicopterBodyHeight/2.0-2,
		h.x+config.HelicopterRotorLen/2.0,
		h.y-config.HelicopterBodyHeight/2.0-2,
		1.0,
		config.ColourMagenta,
		false,
	)
}

func (g *Game) drawHelicopters(screen *ebiten.Image) {
	for _, h := range g.helicopters {
		g.drawHelicopter(screen, h)
	}
}

func (g *Game) spawnHelicopter() {
	if rand.Float32() < config.HelicopterSpawnChance {
		startX := -float32(
			config.HelicopterBodyWidth + config.HelicopterTailWidth,
		)
		startY := float32(50 + rand.Intn(50))
		leftToRight := true
		if rand.Intn(2) == 1 {
			startX = config.ScreenWidth - startX
			leftToRight = false
		}
		g.helicopters = append(g.helicopters, &Helicopter{
			x:           startX,
			y:           startY,
			leftToRight: leftToRight,
			lastDrop:    time.Now(),
		})
	}
}

func (g *Game) updateHelicopters() {
	active := make([]*Helicopter, 0, len(g.helicopters))
	for _, h := range g.helicopters {
		vx := float32(config.HelicopterSpeed)
		if !h.leftToRight {
			vx = -vx
		}
		h.x += vx
		timePassed := time.Since(h.lastDrop)
		if timePassed > config.HelicopterDropRate*time.Second &&
			g.canDrop(h.x) && rand.Float32() < config.ParatrooperSpawnChance {
			g.spawnParatrooper(h.x, h.y)
			h.lastDrop = time.Now()
		}
		if h.x > -100 && h.x < config.ScreenWidth+100 {
			active = append(active, h)
		}
	}
	g.helicopters = active
}

func (g *Game) canDrop(x float32) bool {
	if x < config.ParatrooperWidth/2.0 ||
		x > config.ScreenWidth-config.ParatrooperWidth/2.0 {
		return false
	}
	baseX := (config.ScreenWidth - config.BaseWidth) / 2.0
	pX := x - config.ParatrooperWidth/2.0
	if utils.Overlap1D(pX, config.ParatrooperWidth, baseX, config.BaseWidth) {
		return false
	}
	for _, p := range g.paratroopers {
		if utils.Overlap1D(
			pX-1.0,
			config.ParatrooperWidth+2.0,
			p.x-config.ParatrooperWidth/2.0-1.0,
			config.ParatrooperWidth+2.0,
		) {
			return false
		}
	}
	return true
}
