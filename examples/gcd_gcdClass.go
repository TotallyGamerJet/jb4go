//Expected output for gcdClass.class

package gcd

type gcd_gcdClass struct {
	*java_lang_Object
}

func (arg0 *gcd_gcdClass) init__V() {
	arg0.java_lang_Object = new(java_lang_Object)
	arg0.java_lang_Object.init__V()
}

func gcd_gcdClass_main_RG_V(arg0 []*java_lang_String) {
	var local1 = new(java_util_Scanner)
	local1.init_java_io_InputStream_V(java_lang_System_In)
	java_lang_System_out.println_G_V(New_string_G("Jarrett Kuklis - Assignment #7\n"))
	for {
		java_lang_System_out.print_G_V(New_string_G("Enter first number (-1 to quit):"))
		var local2 = local1.nextInt__I()
		if local2 < 0 {
			break
		}
		java_lang_System_out.print_G_V(New_string_G("Enter second number: "))
		var local3 = local1.nextInt__I()
		var local4 = gcd_gcdClass_gcd_II_I(local2, local3)
		var local5 = new(java_lang_StringBuilder)
		local5.init__V()
		local5 = local5.append_G_java_lang_StringBuilder(New_string_G("GCD is: "))
		local5 = local5.append_I_java_lang_StringBuilder(local4)
		java_lang_System_out.println_G_V(local5.String__G())
	}
	return
}

func gcd_gcdClass_gcd_II_I(arg0 int32, arg1 int32) int32 {
	if arg1 == 0 {
		return arg0
	}
	return gcd_gcdClass_gcd_II_I(arg1, arg0%arg1)
}
