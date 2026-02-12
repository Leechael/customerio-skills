Feature: Version command

  Scenario: Running cio version shows version info
    When I run "cio version"
    Then the exit code should be 0
    And the output should contain "cio"
    And the output should contain "commit:"
    And the output should contain "built:"
