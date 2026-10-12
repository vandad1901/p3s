Feature: Authn Login
    Background:
        Given user registers with the following data
            | Key    | Username      | Email                     | DisplayName | Password       |
            | $User1 | johndoe-{16c} | johndoe-{16c}@example.com | {16c}       | p@ssw0rd-{16c} |
        And user should receive valid JWT with the following data
            | Key    |
            | $User1 |
    Scenario: Login with valid credentials
        When user logs in with the following data
            | Key         | Username          | DisplayName | Password          |
            | $User1Login | {$User1.Username} | {16c}       | {$User1.Password} |
        Then user should receive valid JWT with the following data
            | Key         |
            | $User1Login |
    Scenario: Login with wrong password
        When user logs in with the following data expecting error
            | Key         | Username          | DisplayName | Password          |
            | $User1Login | {$User1.Username} | {16c}       | wrong-password-01 |

        Then user should get the following error
            """
            authn.invalidAuth
            """
    Scenario: Login with unknown username
        When user logs in with the following data expecting error
            | Key         | Username     | DisplayName | Password          |
            | $User1Login | nobody-{16c} | {16c}       | {$User1.Password} |

        Then user should get the following error
            """
            authn.invalidAuth
            """
