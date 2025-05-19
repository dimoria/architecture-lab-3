package lang

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"strconv"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

type Parser struct{}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	var ops []painter.Operation

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		cmd := fields[0]
		args := fields[1:]

		switch cmd {
		case "white":
			ops = append(ops, painter.WhiteFill)
		case "green":
			ops = append(ops, painter.GreenFill)
		case "update":
			ops = append(ops, painter.UpdateOp)
		case "bgrect":
			if len(args) != 4 {
				return nil, fmt.Errorf("bgrect requires 4 arguments")
			}
			rect, err := parseRect(args)
			if err != nil {
				return nil, err
			}
			ops = append(ops, painter.BgRect(rect))
		case "figure":
			if len(args) != 2 {
				return nil, fmt.Errorf("figure requires 2 arguments")
			}
			pos, err := parsePoint(args)
			if err != nil {
				return nil, err
			}
			ops = append(ops, painter.Figure(pos))
		case "move":
			if len(args) != 2 {
				return nil, fmt.Errorf("move requires 2 arguments")
			}
			delta, err := parsePoint(args)
			if err != nil {
				return nil, err
			}
			ops = append(ops, painter.Move(delta))
		case "reset":
			ops = append(ops, painter.ResetOp)
		default:
			return nil, fmt.Errorf("unknown command: %s", cmd)
		}
	}

	return ops, nil
}

func parseRect(args []string) (image.Rectangle, error) {
	coords := make([]float64, 4)
	for i, arg := range args {
		val, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return image.Rectangle{}, err
		}
		coords[i] = val
	}

	return image.Rect(
		int(coords[0]*800),
		int(coords[1]*800),
		int(coords[2]*800),
		int(coords[3]*800),
	), nil
}

func parsePoint(args []string) (image.Point, error) {
	x, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return image.Point{}, err
	}
	y, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return image.Point{}, err
	}
	return image.Point{
		X: int(x * 800),
		Y: int(y * 800),
	}, nil
}
