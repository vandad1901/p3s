Feature: Posts Get
    Scenario: User sees own posts
        Given user logs in as the following user
            | Key    | UserID |
            | $User1 | 1      |
        Given user defines post with the following data
            | Key    | Title         | Slug          | Status |
            | $Post1 | My First Post | my-blog-{64d} | $draft |
        And user adds post blocks with the following data
            | Key          | Position | HeaderKey | BlockType | TextContent                 | Metadata       |
            | $Post1Block1 | 10       | $Post1    | $text     | This is a text block.       | {"size": "h1"} |
            | $Post1Block2 | 20       | $Post1    | $text     | This is another text block. | {"size": "p"}  |
        When user creates post with the following data
            | Key    |
            | $Post1 |

        Given user defines post with the following data
            | Key    | Title         | Slug          | Status     |
            | $Post2 | My First Post | my-blog-{64d} | $published |
        And user adds post blocks with the following data
            | Key          | Position | HeaderKey | BlockType | TextContent                 | Metadata       |
            | $post2Block1 | 10       | $Post2    | $text     | This is a text block.       | {"size": "h1"} |
            | $Post2Block2 | 20       | $Post2    | $text     | This is another text block. | {"size": "p"}  |
        When user creates post with the following data
            | Key    |
            | $Post2 |

        Then user should be able to see posts with the following data
            | Key     | AuthorID | ExpectedPosts |
            | $query1 | 1        | $Post1,$Post2 |
    Scenario: others don't see drafted posts
        Given user logs in as the following user
            | Key    | UserID |
            | $User1 | 1      |
        And user defines post with the following data
            | Key    | Title         | Slug          | Status |
            | $Post1 | My First Post | my-blog-{64d} | $draft |
        And user adds post blocks with the following data
            | Key          | Position | HeaderKey | BlockType | TextContent                 | Metadata       |
            | $Post1Block1 | 10       | $Post1    | $text     | This is a text block.       | {"size": "h1"} |
            | $Post1Block2 | 20       | $Post1    | $text     | This is another text block. | {"size": "p"}  |
        And user creates post with the following data
            | Key    |
            | $Post1 |

        And user defines post with the following data
            | Key    | Title         | Slug          | Status     |
            | $Post2 | My First Post | my-blog-{64d} | $published |
        And user adds post blocks with the following data
            | Key          | Position | HeaderKey | BlockType | TextContent                 | Metadata       |
            | $post2Block1 | 10       | $Post2    | $text     | This is a text block.       | {"size": "h1"} |
            | $Post2Block2 | 20       | $Post2    | $text     | This is another text block. | {"size": "p"}  |
        And user creates post with the following data
            | Key    |
            | $Post2 |

        And user logs in as the following user
            | Key    | UserID |
            | $User2 | 2      |

        Then user should be able to see posts with the following data
            | Key     | AuthorID | ExpectedPosts | MissingPosts |
            | $query1 | 1        | $Post2        | $Post1       |
