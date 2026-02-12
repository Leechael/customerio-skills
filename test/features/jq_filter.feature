Feature: JQ filter

  Background:
    Given a mock API server is running

  Scenario: Filter with --jq
    Given the mock server responds to "GET /v1/segments" with:
      """
      {"segments": [{"id": 1, "name": "Active Users"}, {"id": 2, "name": "Churned"}]}
      """
    When I run "cio segments ls --jq .segments[0].name" against the mock server
    Then the exit code should be 0
    And the output should contain "Active Users"

  Scenario: Invalid jq expression
    Given the mock server responds to "GET /v1/segments" with:
      """
      {"segments": []}
      """
    When I run "cio segments ls --jq .[invalid" against the mock server
    Then the exit code should be 1
