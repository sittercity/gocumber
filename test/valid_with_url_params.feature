Feature: gocumber example feature
  As a person wanting to learn about gocumber
  I want to be shown a bunch of awesome examples

  Scenario: Create a user with a json payload
    When I get "/something/%{UUID}"
