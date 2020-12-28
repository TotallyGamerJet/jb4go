
/*
 * Christian Kuklis AverageNums
 * 
 * Prompts the user for three integer values then prints their average to the
 * nearest tenth.
 * 
 */

import java.util.Scanner;
import java.text.DecimalFormat;

public class AverageNums
{
    public static void main()
    {
        Scanner input = new Scanner(System.in);
        DecimalFormat one = new DecimalFormat("0.0");
        
        int first, second, third;
        double avg;
        
        System.out.print("Enter the 1st number: ");
        first = input.nextInt();
        
        System.out.print("Enter the 2nd number: ");
        second = input.nextInt();
        
        System.out.print("Enter the 3rd number: ");
        third = input.nextInt();
        
        avg = (first + second + third) / 3.0;
        
        System.out.print("\nThe average is " + one.format(avg)); 
        
    }
}
