Feature: Help output

  Scenario: Running cio with no arguments shows help
    When I run "cio"
    Then the exit code should be 0
    And the output should contain "CLI for Customer.io App API"

  Scenario: Running cio --help shows help
    When I run "cio --help"
    Then the exit code should be 0
    And the output should contain "Usage:"
    And the output should contain "Available Commands:"

  Scenario: Running cio segments --help shows subcommand help
    When I run "cio segments --help"
    Then the exit code should be 0
    And the output should contain "segments"
