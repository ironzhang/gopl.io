package tempconv

import "testing"

func TestTempconv(t *testing.T) {
	c := Celsius(BoilingC)
	f := Fahrenheit(0)
	t.Logf("Celsius: %s, Fahrenheit: %s\n", c, CToF(c))
	t.Logf("Fahrenheit: %s, Celsius: %s\n", f, FToC(f))
}
