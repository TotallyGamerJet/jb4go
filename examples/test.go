package main

func init() {
	var currentLabel = 0
controlFlowLoop:
	for {
		switch currentLabel {
		case 0:
			var v0 = @local0
			_ = v0.<init>() //[java/lang/Object.<init>:()V]
			return
			fallthrough
		default:
		}
	}
}

func main() {
	var currentLabel = 0
controlFlowLoop:
	for {
		switch currentLabel {
		case 0:
			var v0 = new(java/util/Scanner)
			var v1 = v0
			var v2 = java/lang/System.in:Ljava/io/InputStream;
			_ = v1.<init>(v2,) //[java/util/Scanner.<init>:(Ljava/io/InputStream;)V]
			@local1 = v0
			var v4 = java/lang/System.out:Ljava/io/PrintStream;
			var v5 = `Jarrett Kuklis - Assignment #7
`
			_ = v4.println(v5,) //[java/io/PrintStream.println:(Ljava/lang/String;)V]
			fallthrough
		case 19:
			var v7 = java/lang/System.out:Ljava/io/PrintStream;
			var v8 = `Enter first number (-1 to quit): `
			_ = v7.print(v8,) //[java/io/PrintStream.print:(Ljava/lang/String;)V]
			var v10 = @local1
			var v11 = v10.nextInt() //[java/util/Scanner.nextInt:()I]
			@local2 = v11
			var v12 = @local2
			if v12 >= 0 { currentLabel = 39; continue controlFlowLoop }
			fallthrough
		case 36:
			currentLabel = 88; continue controlFlowLoop
			fallthrough
		case 39:
			var v13 = java/lang/System.out:Ljava/io/PrintStream;
			var v14 = `Enter second number: `
			_ = v13.print(v14,) //[java/io/PrintStream.print:(Ljava/lang/String;)V]
			var v16 = @local1
			var v17 = v16.nextInt() //[java/util/Scanner.nextInt:()I]
			@local3 = v17
			var v18 = @local2
			var v19 = @local3
			var v20 = gcd/gcdClass.gcd(v19,v18,) //[gcd/gcdClass.gcd:(II)I]
			@local4 = v20
			var v21 = java/lang/System.out:Ljava/io/PrintStream;
			var v22 = new(java/lang/StringBuilder)
			var v23 = v22
			_ = v23.<init>() //[java/lang/StringBuilder.<init>:()V]
			var v25 = `GCD is: `
			var v26 = v22.append(v25,) //[java/lang/StringBuilder.append:(Ljava/lang/String;)Ljava/lang/StringBuilder;]
			var v27 = @local4
			var v28 = v26.append(v27,) //[java/lang/StringBuilder.append:(I)Ljava/lang/StringBuilder;]
			var v29 = v28.toString() //[java/lang/StringBuilder.toString:()Ljava/lang/String;]
			_ = v21.print(v29,) //[java/io/PrintStream.print:(Ljava/lang/String;)V]
			currentLabel = 19; continue controlFlowLoop
			fallthrough
		case 88:
			return
			fallthrough
		default:
		}
	}
}

func gcd() {
	var currentLabel = 0
controlFlowLoop:
	for {
		switch currentLabel {
		case 0:
			var v0 = @local1
			if v0 != 0 { currentLabel = 6; continue controlFlowLoop }
			fallthrough
		case 4:
			var v1 = @local0
			return v1
			fallthrough
		case 6:
			var v2 = @local1
			var v3 = @local0
			var v4 = @local1
			var v5 = v3 % v4
			var v6 = gcd/gcdClass.gcd(v5,v2,) //[gcd/gcdClass.gcd:(II)I]
			return v6
			fallthrough
		default:
		}
	}
}
