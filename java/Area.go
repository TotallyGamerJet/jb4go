package java

type Area struct {
	*java_lang_Object
}

func (arg0 *Area) init__V() {
	var v0 *Area
	v0 = arg0
	v0.java_lang_Object.init__V()
	return
}

func Area_calcSquare_D_D(arg0 float64) float64 {
	var v0 float64
	var v1 float64
	var v2 float64
	v0 = arg0
	v1 = arg0
	v2 = v0 * v1
	return v2
}

func Area_calcRectangle_DD_D(arg0 float64, arg2 float64) float64 {
	var v0 float64
	var v1 float64
	var v2 float64
	v0 = arg0
	v1 = arg2
	v2 = v0 * v1
	return v2
}

func Area_calcTriangle_DD_D(arg0 float64, arg2 float64) float64 {
	var v0 float64
	var v1 float64
	var v2 float64
	var v3 float64
	var v4 float64
	v0 = arg0
	v1 = arg2
	v2 = v0 * v1
	v3 = 2e+00
	v4 = v2 / v3
	return v4
}

func Area_calcTrapezoid_DDD_D(arg0 float64, arg2 float64, arg4 float64) float64 {
	var v0 float64
	var v1 float64
	var v2 float64
	var v3 float64
	var v4 float64
	var v5 float64
	var v6 float64
	v0 = arg0
	v1 = arg2
	v2 = v0 + v1
	v3 = arg4
	v4 = v2 * v3
	v5 = 2e+00
	v6 = v4 / v5
	return v6
}

func Area_calcCircle_D_D(arg0 float64) float64 {
	var v0 float64
	var v1 float64
	var v2 float64
	var v3 float64
	var v4 float64
	v0 = 3.141592653589793e+00
	v1 = arg0
	v2 = 3e+00
	v3 = java_lang_Math_pow_DD_D(v1, v2)
	v4 = v0 * v3
	return v4
}
