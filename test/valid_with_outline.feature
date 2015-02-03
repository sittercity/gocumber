Feature: Example Feature with Outline
  As a user yadda yadda
  I want to do stuff

  Scenario Outline: Create user with bad data
    Given I have no users
    When I create a new user with the following data:
      | key      | value      |
      | username | <username> |
      | first    | <first>    |
      | last     | <last>     |
    And no users should be created

    Examples:
      | username | first | last  |
      |          |       |       |
      |          | fname | lname |
      | newuser  |       | lname |
      | newuser  | fname |       |
      | newuser  | fname | lname |
