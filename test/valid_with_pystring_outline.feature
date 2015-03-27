Feature: Scenario Outline with DocString

Scenario Outline: Something
  Given something is:
    """
    <key> functional
    """

Examples:
  | key       |
  | minimally |
