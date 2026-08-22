Feature: Post Delete
    Background:
        Given user defines post with the following data
            | Key    | Title         | Slug          |
            | $Post1 | My First Post | my-blog-{64d} |
        And user adds post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent                 | Metadata        |
            | $Block1 | 10       | $Post1    | $text     | This is a text block.       | {"header": "3"} |
            | $Block2 | 20       | $Post1    | $text     | This is another text block. | {"header": "3"} |
        And user creates post with the following data
            | Key    |
            | $Post1 |
    Scenario: Delete a post
        When user deletes post with the following data
            | Key    |
            | $Post1 |

        Then user should not be able to see post with the following data
            | Key    |
            | $Post1 |
    Scenario: Delete a post with invalid UpdatedAt
        When user deletes post with the following data expecting error
            | Key    | UpdatedAt |
            | $Post1 | {$empty}  |

        Then user should get the following error
            """
            post.NotFound
            """
