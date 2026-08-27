Feature: Authn Refresh JWT
    Background:
        Given user registers with the following data
            | Key    | Username      | Email                     | Password       |
            | $User1 | johndoe-{16c} | johndoe-{16c}@example.com | p@ssw0rd-{16c} |
        And user should receive valid JWT with the following data
            | Key    |
            | $User1 |
    Scenario: Refresh JWT with a valid session
        When user refreshes jwt with the following data
            | Key    |
            | $User1 |
        Then user should receive valid JWT with the following data
            | Key    |
            | $User1 |
    Scenario: Refresh JWT with an invalid refresh token
        When user refreshes jwt with the following data expecting error
            | Key    | RefreshToken       |
            | $User1 | invalid-token-1234 |

        Then user should get the following error
            """
            authn.invalidAuth
            """
    Scenario: Refresh JWT with an unknown session
        When user refreshes jwt with the following data expecting error
            | Key    | SessionID |
            | $User1 | 999999    |

        Then user should get the following error
            """
            authn.invalidAuth
            """
    Scenario: Refresh JWT with an unknown user
        When user refreshes jwt with the following data expecting error
            | Key    | UserID |
            | $User1 | 999999 |

        Then user should get the following error
            """
            authn.invalidAuth
            """
