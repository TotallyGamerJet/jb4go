package main

/* TYPES:
JAVA : GO
char : uint16
short : int16
byte : int8
int : int32
long : int64
float : float32
double : float64
boolean : bool
*/

/*
Conversion Rules:
byte - B
char - C
double - D
float - F
int - I
long - J
String - L{String}?
short - S
boolean - Z
array - R{type}
void - V
*/

/* Methods and Functions:
All methods start with its name exported followed by an '_' then a list of the parameters following the
conversion rules
than another '_' and its return type. If there are no parameters there are two '__' next to each
other.
If a method is static it becomes a go function with the P_java_lang_Object's name preceeding it.
Constructors are the same as methods expect they have an '_C' at the very end to signify its a
constructor
*/

/*
Each object has two automatic methods Synchronize and Unsynchronize that lock and unlock the mutex
on the P_java_lang_Object.
*/
