# Let's build a Java Virtual Machine

The JVM specification (https://docs.oracle.com/javase/specs/jvms/se17/jvms17.pdf)
is quiet big. What is the smallest increment we want to target first?

In a previous attempt, I considered a simple Hello World example.
But accessing System.out.println() turned out to include a whole baggage of things.
I have learned by lesson though, and will focus on simplicity above else.
