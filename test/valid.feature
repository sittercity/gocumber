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
