Feature: Authn Register
    Scenario: Register a new user
        When user registers with the following data
            | Key    | Username      | Email                     | DisplayName | Password       |
            | $User1 | johndoe-{16c} | johndoe-{16c}@example.com | {16c}       | p@ssw0rd-{16c} |

        Then user should receive valid JWT with the following data
            | Key    |
            | $User1 |
    Scenario: Register with short username
        When user registers with the following data expecting error
            | Key    | Username | Email                     | DisplayName | Password       |
            | $User1 | {4c}     | johndoe-{16c}@example.com | {16c}       | p@ssw0rd-{16c} |

        Then user should get the following error
            """
            validating user: identity.validation.invalidUsername
            """
    Scenario: Register with short display name
        When user registers with the following data expecting error
            | Key    | Username  | Email                     | DisplayName | Password       |
            | $User1 | doe-{16c} | johndoe-{16c}@example.com | {2c}        | p@ssw0rd-{16c} |

        Then user should get the following error
            """
            validating user: identity.validation.invalidDisplayName
            """
    Scenario: Register with empty email
        When user registers with the following data expecting error
            | Key    | Username      | Email    | DisplayName | Password       |
            | $User1 | johndoe-{16c} | {$empty} | {16c}       | p@ssw0rd-{16c} |

        Then user should get the following error
            """
            validating user: identity.validation.emptyEmail
            """
    Scenario: Register with invalid email
        When user registers with the following data expecting error
            | Key    | Username      | Email        | DisplayName | Password       |
            | $User1 | johndoe-{16c} | not-an-email | {16c}       | p@ssw0rd-{16c} |

        Then user should get the following error
            """
            validating user: identity.validation.invalidEmail
            """
    Scenario: Register with insecure password
        When user registers with the following data expecting error
            | Key    | Username      | Email                     | DisplayName | Password |
            | $User1 | johndoe-{16c} | johndoe-{16c}@example.com | {16c}       | short    |

        Then user should get the following error
            """
            validating password: credential.validation.InsecurePassword
            """
    Scenario: Register with duplicate username
        Given user registers with the following data
            | Key    | Username   | Email                  | DisplayName | Password       |
            | $User1 | dupe-{16c} | dupe-{16c}@example.com | {16c}       | p@ssw0rd-{16c} |

        When user registers with the following data expecting error
            | Key    | Username          | Email                   | DisplayName | Password       |
            | $User2 | {$User1.Username} | other-{16c}@example.com | {16c}       | p@ssw0rd-{16c} |

        Then user should get the following error
            """
            creating user: identity.DuplicateUsername
            """
