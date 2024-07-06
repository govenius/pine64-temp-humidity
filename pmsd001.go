package main

import (
	"fmt"
	"golang.org/x/exp/io/i2c"
	"time"
	"flag"
)

const (
	TEMP_NO_HOLD = byte(0xF3)
	HUMIDITY     = byte(0xF5)
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func waitForRead() {
	time.Sleep(500 * time.Millisecond)
}

func getTemp(d i2c.Device) (float64, float64) {
	err := d.Write([]byte{TEMP_NO_HOLD})
	check(err)

	waitForRead()

	response := make([]byte, 2)
	check(d.Read(response))

	temp := (int64(response[0]) * 256 + int64(response[1])) & 0xFFFC
	tempC := -46.85 + (175.72 * float64(temp) / 65536.0)
	tempF := tempC * 1.8 + 32

	return tempF, tempC
}

func getRelativeHumidity(d i2c.Device) float64 {
	err := d.Write([]byte{HUMIDITY})
	check(err)

	waitForRead()

	response := make([]byte, 2)
	check(d.Read(response))

	data := (int64(response[0]) * 256 + int64(response[1])) & 0xFFFC
	humidity := (125.0 * float64(data) / 65536.0) - 6.0

	return humidity
}

func main() {

	// Define and parse comman line arguments
	i2cBus := flag.String("i", "/dev/i2c-1", "i2c bus to query (default: /dev/i2c-1)")
	i2cAddress := flag.Int("a", 0x40, "i2c address to query (default: 0x40)")
	jsonOutput := flag.Bool("j", false, "Output JSON instead of tab separated data (default: false)")
	sensorName := flag.String("n", "pi", "sensor name included in JSON output (default: pi)")
	flag.Parse()

	d, err := i2c.Open(&i2c.Devfs{Dev: *i2cBus}, *i2cAddress)
	check(err)

	tempF, tempC := getTemp(*d)
	humidity := getRelativeHumidity(*d)

	if *jsonOutput {
		fmt.Printf("{\"%s\": ", *sensorName)
		fmt.Printf("{\"time\": %d, \"temperature\": %.2f , \"humidity\": %.2f}", time.Now().Unix(), tempC, humidity)
		fmt.Printf("}\n")
	} else {
		fmt.Printf("%d\t%.2f\t%.2f\t%.2f\n", time.Now().Unix(), tempC, tempF, humidity)
	}
}
