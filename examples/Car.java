/**
  This class contains the data object for the car object

  @author Christian Kuklis
  @version 4/15/2020
*/

public class Car
{
  private int year, speed;
  private String make;

  public Car(int year, String make)
  {
    this.year = year;
    this.make = make;
  }

  public int getYear()
  {
    return year;
  }

  public String getMake()
  {
    return make;
  }

  public int getSpeed()
  {
    return speed;
  }

  public void accelerate()
  {
    speed += 5;
  }

  public void brake()
  {
    speed -= 3;
  }  
}
