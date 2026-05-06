package complexnumbers

import "math"

// Define the Number type here.
type Number struct {
    real float64
    imag float64
}
func (n Number) Real() float64 {
	return n.real
}

func (n Number) Imaginary() float64 {
	return n.imag
}

func (n1 Number) Add(n2 Number) Number {
	return Number{
        real: n1.real + n2.real,
        imag: n1.imag + n2.imag,
    }
}

func (n1 Number) Subtract(n2 Number) Number {
	return Number{
        real: n1.real - n2.real,
        imag: n1.imag - n2.imag,
    }
}

func (n1 Number) Multiply(n2 Number) Number {
	return Number{
        real: (n1.real * n2.real) - (n1.imag * n2.imag),
		imag: (n1.real * n2.imag) + (n1.imag * n2.real),
    }
}

func (n Number) Times(factor float64) Number {
	return Number{
        real: n.real * factor,
        imag: n.imag * factor,
    }
}

func (n1 Number) Divide(n2 Number) Number {
	divisor := (n2.real * n2.real) + (n2.imag * n2.imag)
    return Number{
        real: ((n1.real * n2.real) + (n1.imag * n2.imag)) / divisor,
		imag: ((n1.imag * n2.real) - (n1.real * n2.imag)) / divisor,
    }
}

func (n Number) Conjugate() Number {
	return Number{
        real: n.real,
        imag: -n.imag,
    }
}

func (n Number) Abs() float64 {
	return math.Hypot(n.real, n.imag)
}

func (n Number) Exp() Number {
	exp := math.Exp(n.real)
    return Number{
        real: exp * math.Cos(n.imag),
        imag: exp * math.Sin(n.imag),
    }
}
