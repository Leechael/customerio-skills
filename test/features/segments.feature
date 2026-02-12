Feature: Segments commands with mock server

  Background:
    Given a mock API server is running

  Scenario: List segments
    Given the mock server responds to "GET /v1/segments" with:
      """
      {"segments": [{"id": 1, "name": "Active Users"}]}
      """
    When I run "cio segments ls" against the mock server
    Then the exit code should be 0
    And the output should contain "Active Users"

  Scenario: Get a segment
    Given the mock server responds to "GET /v1/segments/42" with:
      """
      {"segment": {"id": 42, "name": "Power Users"}}
      """
    When I run "cio segments get 42" against the mock server
    Then the exit code should be 0
    And the output should contain "Power Users"

  Scenario: Missing API token
    When I run "cio segments ls" without an API token
    Then the exit code should be 1
    And the output should contain "CUSTOMERIO_API_TOKEN"
