package gcd;

import java.util.Scanner;

public class gcdClass {

    public static void main(String[] args) {
        Scanner s = new Scanner(System.in);
        System.out.println("Jarrett Kuklis - Assignment #7\n");
        while (true) {
            System.out.print("Enter first number (-1 to quit): ");
            int f = s.nextInt();
            if (f < 0) {
                break;
            }
            System.out.print("Enter second number: ");
            int sec = s.nextInt();
            int gcd = gcd(f, sec);
            System.out.println("GCD is: "+ gcd);
        }
    }

    public static int gcd(int x, int y) {
        if (y == 0) {
            return x;
        }
        return gcd(y, x%y);
    }
}