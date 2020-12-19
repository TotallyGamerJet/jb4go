package employee;

public class Employee {
	private final String firstName;
    private final String lastName;
    private final String socialSecurityNumber;

    public Employee(String firstName, String lastName, String socialSecurityNumber) {
	      this.firstName = firstName;
	      this.lastName = lastName;
	      this.socialSecurityNumber = socialSecurityNumber;
    }

    // return first name
    public String getFirstName() {
       return firstName;
    }

    // return last name
    public String getLastName() {
      return lastName;
    }

    // return social security number
    public String getSocialSecurityNumber() {
      return socialSecurityNumber;
    }

}