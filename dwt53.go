package dwt

/**
 *  dwt53.c - Fast discrete biorthogonal CDF 5/3 wavelet forward and inverse transform (lifting implementation)
 *
 *  This code is provided "as is" and is given for educational purposes.
 *  2007 - Gregoire Pau - gregoire.pau@ebi.ac.uk
 */

const (
	p1_53  = -0.5
	ip1_53 = -p1_53

	u1_53  = 0.25
	iu1_53 = -u1_53

	scale53  = 1.4142135623730951 // math.Sqrt(2.0)
	iscale53 = 1 / scale53
)

// Fwt53 performs a bi-orthogonal 5/3 wavelet transformation (lifting implementation)
// of the signal in slice xn. The length of the signal n = len(xn) must be a power of 2.
//
// The input in slice xn will be replaced by the transformation:
//
// The first half part of the output signal contains the approximation coefficients.
// The second half part contains the detail coefficients (aka. the wavelets coefficients).
func Fwt53(xn []float64) {
	n := validateLen(xn)

	// predict 1
	for i := 1; i < n-2; i += 2 {
		xn[i] += p1_53 * (xn[i-1] + xn[i+1])
	}
	xn[n-1] += 2 * p1_53 * xn[n-2]

	// update 1
	for i := 2; i < n; i += 2 {
		xn[i] += u1_53 * (xn[i-1] + xn[i+1])
	}
	xn[0] += 2 * u1_53 * xn[1]

	// scale
	for i := 0; i < n; i++ {
		if i%2 != 0 {
			xn[i] = xn[i] * scale53
		} else {
			xn[i] /= scale53
		}
	}

	// pack
	tb := make([]float64, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			tb[i/2] = xn[i]
		} else {
			tb[n/2+i/2] = xn[i]
		}
	}
	copy(xn, tb)
}

// Iwt53 performs an inverse bi-orthogonal 5/3 wavelet transformation of xn.
// This is the inverse function of Fwt53 so that Iwt53(Fwt53(xn))=xn for every signal xn of length n.
//
// The length of slice xn must be a power of 2.
//
// The coefficients provided in slice xn are replaced by the original signal.
func Iwt53(xn []float64) {
	n := validateLen(xn)

	// unpack
	tb := make([]float64, n)
	for i := 0; i < n/2; i++ {
		tb[i*2] = xn[i]
		tb[i*2+1] = xn[i+n/2]
	}
	copy(xn, tb)

	// undo scale
	for i := 0; i < n; i++ {
		if i%2 != 0 {
			xn[i] *= iscale53
		} else {
			xn[i] /= iscale53
		}
	}

	// undo update 1
	for i := 2; i < n; i += 2 {
		xn[i] += iu1_53 * (xn[i-1] + xn[i+1])
	}
	xn[0] += 2 * iu1_53 * xn[1]

	// undo predict 1
	for i := 1; i < n-2; i += 2 {
		xn[i] += ip1_53 * (xn[i-1] + xn[i+1])
	}
	xn[n-1] += 2 * ip1_53 * xn[n-2]
}
