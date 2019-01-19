package termui

import (
	"image"
)

type Canvas struct {
	CellMap map[image.Point]Cell
	Block
}

func NewCanvas() *Canvas {
	return &Canvas{
		Block:   *NewBlock(),
		CellMap: make(map[image.Point]Cell),
	}
}

// given points correspond to dots within a braille character, not cells
func (self *Canvas) Line(p0, p1 image.Point, attr Attribute) {
	leftPoint, rightPoint := p0, p1
	if leftPoint.X > rightPoint.X {
		leftPoint, rightPoint = rightPoint, leftPoint
	}

	xDistance := AbsInt(leftPoint.X - rightPoint.X)
	yDistance := AbsInt(leftPoint.Y - rightPoint.Y)
	slope := float64(yDistance) / float64(xDistance)
	slopeDirection := 1
	if rightPoint.Y < leftPoint.Y {
		slopeDirection = -1
	}

	targetYCoordinate := float64(leftPoint.Y)
	currentYCoordinate := leftPoint.Y
	for i := leftPoint.X; i < rightPoint.X; i++ {
		targetYCoordinate += (slope * float64(slopeDirection))
		if currentYCoordinate == int(targetYCoordinate) {
			point := image.Pt(i/2, currentYCoordinate/4)
			self.CellMap[point] = Cell{
				self.CellMap[point].Rune | BRAILLE[currentYCoordinate%4][i%2],
				AttrPair{attr, ColorDefault},
			}
		}
		for currentYCoordinate != int(targetYCoordinate) {
			point := image.Pt(i/2, currentYCoordinate/4)
			self.CellMap[point] = Cell{
				self.CellMap[point].Rune | BRAILLE[currentYCoordinate%4][i%2],
				AttrPair{attr, ColorDefault},
			}
			currentYCoordinate += slopeDirection
		}
	}
}

func (self *Canvas) Draw(buf *Buffer) {
	for point, cell := range self.CellMap {
		if point.In(self.Rectangle) {
			buf.SetCell(Cell{cell.Rune + BRAILLE_OFFSET, cell.Attrs}, point)
		}
	}
}
