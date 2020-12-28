import java.util.Scanner;

/**
  This program tests the functionality of the Car object class

  @author Christian Kuklis
  @version 4/15/2020
*/

public class CarDemo
{
  public static void main(String[] args)
  {
    int year;
    String make;

    Scanner input = new Scanner(System.in);

    System.out.print("Enter the year of the car: ");
    year = input.nextInt();

    input.nextLine();

    System.out.print("Enter the make of the car: ");
    make = input.nextLine();

    Car car = new Car(year, make);

    System.out.println();
    System.out.println("The year of the car is " + car.getYear() + " and the make is " + car.getMake());
    System.out.println();

    for(int i = 0; i < 10; i++)
    {
      car.accelerate();
      System.out.println("You are now travelling " + car.getSpeed() + " mph");
    }

    for(int i = 0; i < 4; i++)
    {
      car.brake();
      System.out.println("You are now travelling " + car.getSpeed() + " mph");
    }
  }
}
