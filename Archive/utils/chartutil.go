package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wcharczuk/go-chart/v2"
)

type Points struct {
	AVG []float64
	P50 []float64
	P90 []float64
	P99 []float64
}

func NewPoints() *Points {
	return &Points{[]float64{}, []float64{}, []float64{}, []float64{}}
}

func (points *Points) LoadDataFromCVS(path string, name string, label string) error {
	fullpath := filepath.Join(path, name)
	csvfile, err := os.Open(fullpath)
	if err != nil {
		return err
	}
	defer csvfile.Close()
	//
	scan := bufio.NewScanner(csvfile)
	for scan.Scan() {
		line := scan.Text()
		item := strings.Split(line, ",")
		if item[0] == label {
			f, err := strconv.ParseFloat(item[3], 64)
			if err != nil {
				continue
			}
			points.AVG = append(points.AVG, f)
			f, err = strconv.ParseFloat(item[4], 64)
			if err != nil {
				continue
			}
			points.P50 = append(points.P50, f)
			f, err = strconv.ParseFloat(item[5], 64)
			if err != nil {
				continue
			}
			points.P90 = append(points.P90, f)
		}
	}
	//
	return nil
}

func (points Points) RenderChart(path string, name string) error {
	graph := chart.Chart{
		Title: strings.TrimSuffix(name, filepath.Ext(name)),
		Background: chart.Style{
			Padding: chart.Box{
				Top:    25,
				Left:   25,
				Right:  25,
				Bottom: 25,
			},
		},

		XAxis: chart.XAxis{
			Name: "Collections",
			ValueFormatter: func(v interface{}) string {
				if vf, isFloat := v.(float64); isFloat {
					return fmt.Sprintf("%.1f", vf)
				}
				return ""
			},
		},

		YAxis: chart.YAxis{
			Name: "TimeCost(ms)",
			ValueFormatter: func(v interface{}) string {
				if vf, isFloat := v.(float64); isFloat {
					return fmt.Sprintf("%.1f", vf)
				}
				return ""
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "AVG Line",
				XValues: chart.Seq{Sequence: chart.NewLinearSequence().WithStart(1).WithEnd(float64(len(points.AVG)))}.Values(),
				YValues: points.AVG,
			},
			chart.ContinuousSeries{
				Name:    "P50 Line",
				XValues: chart.Seq{Sequence: chart.NewLinearSequence().WithStart(1).WithEnd(float64(len(points.P50)))}.Values(),
				YValues: points.P50,
			},
			chart.ContinuousSeries{
				Name:    "P90 Line",
				XValues: chart.Seq{Sequence: chart.NewLinearSequence().WithStart(1).WithEnd(float64(len(points.P90)))}.Values(),
				YValues: points.P90,
			},
		},
	}
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}
	//
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		return err
	}
	//
	fullpath := filepath.Join(path, name)
	chartfile, err := os.OpenFile(fullpath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer chartfile.Close()
	writer := bufio.NewWriter(chartfile)
	writer.Write(buffer.Bytes())
	err = writer.Flush()
	if err != nil {
		return err
	}
	//
	return nil
}
