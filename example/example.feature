Feature: gocumber example feature
  As a person wanting to learn about gocumber
  I want to be shown a bunch of awesome examples

  Scenario: Create a user with a json payload
    When I create a user with the following json data:
    """
    {
      "username": "CoolDude",
      "first": "Cool",
      "last": "Dude",
      "activated": true
    }
    """
    Then the user should be created with the expected data

  Scenario: Create a user with table data
    When I create a user with the following table data:
      | key       | value    |
      | username  | CoolDude |
      | first     | Cool     |
      | last      | Dude     |
      | activated | true     |
    Then the user should be created with the expected data

  Scenario: Deactivate a user
    Given I have an existing user with the name "Cool Dude"
    When I deactivate the user named "Cool Dude"
    Then the user "Cool Dude" should be deactivated

  Scenario: Deactivate a user using more generic given step
    Given I have an existing user
    When I deactivate the user named "Cool Dude"
    Then the user "Cool Dude" should be deactivated

  Scenario: Read a user and check against table data
    Given I create a user with the following json data:
    """
    {
      "username": "CoolDude",
      "first": "Cool",
      "last": "Dude",
      "activated": true
    }
    """
    When I read the user named "CoolDude"
    Then I should see the following data:
      | username | first | last |
      | CoolDude | Cool  | Dude |

  Scenario Outline: Create user with bad data
    Given I have no users
    When I create a user with the following table data:
      | key      | value      |
      | username | <username> |
      | first    | <first>    |
      | last     | <last>     |
    Then no users should be created

    Examples:
      | username | first | last  |
      |          |       |       |
      |          | fname | lname |
      | uuid     |       | lname |
      | uuid     | fname |       |
      | uuid     | fname | lname |
