package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"

	"cloud.google.com/go/civil"
	"github.com/briandowns/spinner"
)

func SpinnerMessage(message string, duration time.Duration) {
	s := spinner.New(spinner.CharSets[4], 100*time.Millisecond)
	s.Suffix = message
	s.Start()
	time.Sleep(duration)
	s.Stop()
}

func LineBreak(num int) {
	line := 1
	for line <= num {
		fmt.Println()
		line += 1
	}
}

func DateToString(d *civil.Date) string {
	return fmt.Sprintf("%d.%d.%d", d.Day, d.Month, d.Year)
}

func FloatToMoney(val float64) float64 {

	output := math.Pow(10, float64(2))
	round := val * output
	return float64(int(round+math.Copysign(0.5, round))) / output
}

func DiscountedPrice(price float64, discount float64) float64 {
	return FloatToMoney(price - ((discount / 100) * price))
}

func PressEnterToContinue() {
	fmt.Println()
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	fmt.Println()
}

func ShowUserInSystem(message string) {
	LineBreak(2)
	fmt.Println("-----------------------------------------")
	fmt.Println(message)
	fmt.Println("-----------------------------------------")
	LineBreak(1)
	PressEnterToContinue()
}

func LoadCSV(name string) string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	filePath := fmt.Sprintf("%s/%s/%s/%s", path, "backend", "db", name)
	_, err = os.Stat(filePath)
	if err != nil {
		panic(fmt.Errorf("no file exists in %s", filePath))
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func WriteJSON(bytes []byte, filename string) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	filePath := fmt.Sprintf("%s/%s/%s/%s", path, "backend", "db", filename)
	_ = ioutil.WriteFile(filePath, bytes, 0644)
}

func LoadJSON(filename string) ([]byte, error) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	filePath := fmt.Sprintf("%s/%s/%s/%s", path, "backend", "db", filename)
	bytes, err := ioutil.ReadFile(filePath)
	return bytes, err
}
