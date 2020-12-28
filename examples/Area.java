/**
* This utility class stores the methods for calculating the area of the various
* geometric shapes.
*
* @author Christian Kuklis
* @version 5/6/2020
*/

public class Area
{
  public static double calcSquare(double side)
  {
    return side * side;
  }

  public static double calcRectangle(double length, double width)
  {
    return length * width;
  }

  public static double calcTriangle(double base, double height)
  {
    return (base * height) / 2.0;
  }

  public static double calcTrapezoid(double base1, double base2, double height)
  {
    return (base1 + base2) * height / 2.0;
  }

  public static double calcCircle(double radius)
  {
    return Math.PI * Math.pow(radius, 3);
  }
}
