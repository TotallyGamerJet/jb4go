//Expected output for gcdClass.class

package gcd

import . "github.com/totallygamerjet/jb4go/java"

type P_gcd_gcdClass struct {
	*P_java_lang_Object
}

func (arg0 *P_gcd_gcdClass) P_init__V() {
	arg0.P_java_lang_Object = new(P_java_lang_Object)
	arg0.P_java_lang_Object.P_init__V()
}

func P_main_RG_V(arg0 []*P_java_lang_String) {
	var local1 = new(P_java_util_Scanner)
	local1.P_init_java_io_InputStream_V(P_java_lang_System_In)
	P_java_lang_System_out.P_println_G_V(NewString_string_G("Jarrett Kuklis - Assignment #7\n"))
	for {
		P_java_lang_System_out.P_print_G_V(NewString_string_G("Enter first number (-1 to quit):"))
		var local2 = local1.P_nextInt__I()
		if local2 < 0 {
			break
		}
		P_java_lang_System_out.P_print_G_V(NewString_string_G("Enter second number: "))
		var local3 = local1.P_nextInt__I()
		var local4 = P_gcd_II_I(local2, local3)
		var local5 = new(P_java_lang_StringBuilder)
		local5.P_init__V()
		local5 = local5.P_append_G_java_lang_StringBuilder(NewString_string_G("GCD is: "))
		local5 = local5.P_append_I_java_lang_StringBuilder(local4)
		P_java_lang_System_out.P_println_G_V(local5.P_String__G())
	}
	return
}

func P_gcd_II_I(arg0 int32, arg1 int32) int32 {
	if arg1 == 0 {
		return arg0
	}
	return P_gcd_II_I(arg1, arg0%arg1)
}
