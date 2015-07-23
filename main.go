//Minesweeper Backend

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

const msg = "Hi this is Allan from PermissionData, calling about TEE-KEE-TEE Blah Blah Blah. I would like to take Web 10 off the load balancer"
const mine = "    X    "
const covered = "    -    "
const uncovered = "         "
const border = "<!------------------Matrix Grid ---------------------->"

type Cell struct {
	MineProximity int
	ID            int
	PosX          int
	PosY          int
	Swept         bool
}

func (c Cell) IsMine() bool {
	return c.MineProximity < 0
}

func (c Cell) PPrint() string {
	return fmt.Sprintf("Id: %v  X: %v  Y: %v IsMine: %v Swept: %v\n", c.ID, c.PosX, c.PosY, c.MineProximity, c.Swept)
}

type Matrix struct {
	Grid      [][]*Cell
	CellCount int
}

func NewMatrix(size int) Matrix {
	var m Matrix
	if size < 3 {
		panic("Grid will be too small")
	}
	m.generateGrid(size)
	return m
}

func (m *Matrix) generateGrid(size int) {
	m.Grid = make([][]*Cell, size)

	for i := 0; i < size; i++ {
		var (
			row = make([]*Cell, size)
		)
		for j := 0; j < len(row); j++ {
			id := m.CellCount
			m.CellCount = m.CellCount + 1
			row[j] = &Cell{ID: id, Swept: false, PosX: i, PosY: j}

			//Assign cell mineProximity to -1 if you want it to have a mine
			if rand.Intn(100) > 65 {
				row[j].MineProximity = -1
			}
		}
		m.Grid[i] = row
	}
}

func (m Matrix) PPrint() {
	fmt.Println(border)
	var symbol string
	for i := 0; i < len(m.Grid); i++ {
		for j := 0; j < len(m.Grid[i]); j++ {
			c := m.Grid[i][j]
			//fmt.Println(c.PPrint())
			if c.Swept {
				if c.IsMine() {
					symbol = covered
				} else if c.MineProximity > 0 {
					symbol = "    " + string(c.MineProximity) + "    "
				} else {
					symbol = uncovered
				}
			} else {
				symbol = covered
			}
			fmt.Print(symbol)
		}
		fmt.Println()
	}

	fmt.Println(border)
}

func (m Matrix) checkDead(c Cell) bool {
	return c.IsMine()
}

func (m Matrix) EndGame() {
	fmt.Println(msg)
	os.Exit(3)
}

func (m Matrix) FindCell(id int) (*Cell, error) {
	var c *Cell

	if !(0 < id && id < m.CellCount) {
		return c, fmt.Errorf("Cannot find cell with id %v\n", id)
	}

	for i := 0; i < len(m.Grid); i++ {
		for j := 0; j < len(m.Grid[i]); j++ {
			cell := m.Grid[i][j]
			if cell.ID == id {
				c = cell
			}
		}
	}

	return c, nil

}

func (m Matrix) CheckCell(c *Cell) {
	if m.checkDead(*c) {
		fmt.Printf("You choose a mine!! %v\n", c.PPrint())
		m.EndGame()
		return
	}

	m.sweepBoard(c)
}
func (m Matrix) sweepBoard(c *Cell) {
	initialSweep := m.sweepCell(c)
	m.sweep(initialSweep)
}

func (m Matrix) sweepCell(c *Cell) []*Cell {
	var n []*Cell
	if !c.Swept {
		c.Swept = true
		n = m.getSafeNeighborsAndRate(c)
	}
	return n

}

func (m Matrix) sweep(targets []*Cell) {
	if len(targets) < 1 {
		return //base case
	}

	//Pop :https://code.google.com/p/go-wiki/wiki/SliceTricks
	t, targets := targets[len(targets)-1], targets[:len(targets)-1]
	targets = append(targets, m.sweepCell(t)...)
	m.sweep(targets)
}

func (m Matrix) getSafeNeighborsAndRate(c *Cell) []*Cell {
	var n []*Cell

	c.MineProximity = 0
	//fmt.Printf("Chosen Cell: %v\n", c.PPrint())

	for i := c.PosX - 1; i < c.PosX+1; i++ {
		//fmt.Printf("i: %v length: %v\n", i, len(m.Grid))

		for j := c.PosY - 1; j < c.PosY+1; j++ {
			//fmt.Printf("j: %v width: %v\n", j, len(m.Grid[0]))

			if i >= 0 && i < len(m.Grid) && j >= 0 && j < len(m.Grid[i]) {

				currNeigh := m.Grid[i][j]
				if currNeigh.IsMine() {
					//fmt.Println("Found a unclicked mine")
					currNeigh.Swept = true // No need to check her again
					c.MineProximity++
				} else if !currNeigh.Swept {
					n = append(n, currNeigh)
				}
			}
		}
	}

	//fmt.Printf("Neighbors found %v", n)

	//Don't Sweep beyond a cell that has mine nieghbors
	fmt.Printf("MineProximity found %v\n", c.MineProximity)
	if c.MineProximity > 0 {
		var empty []*Cell
		return empty
	}

	return n

}

func main() {
	//args := os.Args[1:]
	m := NewMatrix(6)

	var (
		guessCell *Cell
		input     string
	)
	for {
		fmt.Printf("Pick a Cell # between 0 and %v:\n", m.CellCount-1)
		m.PPrint()
		if _, err := fmt.Scanf("%s", &input); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		cellId, err := strconv.Atoi(input)

		if err != nil {
			fmt.Printf("%s\n is not a number\n", err)
		} else if guessCell, err = m.FindCell(cellId); err != nil {
			fmt.Printf(err.Error())
		} else {
			m.CheckCell(guessCell)
		}
	}
}
