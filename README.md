# jb4go
 Java Bytecode to Go Transpiler
 
## FAQ
### Why Java Bytecode and not the Source?
Java bytecode is rigid and well define while java source code uses the compiler to add syntactic sugar.
Also, many languages compile to java bytecode so this transpiler will also work for any language
that can produce java bytecode. 

## Resources
I am not an expert on compilers or transpilers and I only have a basic understanding of java bytecode. I used the following
resources to learn and base this implementation on.
1. https://www.mirkosertic.de/blog/2017/06/compiling-bytecode-to-javascript/
1. https://tomassetti.me/how-to-write-a-transpiler/
1. https://en.wikipedia.org/wiki/Java_bytecode_instruction_listings
1. https://docs.oracle.com/javase/specs/jvms/se7/html/jvms-6.html