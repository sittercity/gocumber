
Feature: Scenario Outline with DocString

Scenario Outline: Something
  Given something is:
    """
    <how> functional <really_though>
    """

Examples:
  | how        | really_though  |
  | minimally  | barely         |
  | incredibly | overwhelmingly |
