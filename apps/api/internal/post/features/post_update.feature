Feature: Post Update
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
    Scenario: Update a post
        When user updates post with the following data
            | Key    | Title           |
            | $Post1 | My Updated Post |

        Then user should be able to see post with the following data
            | Key    |
            | $Post1 |
    Scenario: Update post blocks
        When user adds post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent            | Metadata        |
            | $Block3 | 30       | $Post1    | $text     | This is a third block. | {"header": "4"} |
        And user edits post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent               | Metadata        |
            | $Block1 | 10       | $Post1    | $text     | This is an updated block. | {"header": "2"} |
        And user removes post blocks with the following data
            | Key     |
            | $Block2 |
        And user updates post with the following data
            | Key    | Title           |
            | $Post1 | My Updated Post |

        Then user should be able to see post with the following data
            | Key    |
            | $Post1 |
    Scenario: Update a post with empty title
        When user updates post with the following data expecting error
            | Key    | Title    |
            | $Post1 | {$empty} |

        Then user should get the following error
            """
            post.validation.EmptyTitle
            """
    Scenario: Update a post with empty slug
        When user updates post with the following data expecting error
            | Key    | Slug     |
            | $Post1 | {$empty} |

        Then user should get the following error
            """
            post.validation.EmptySlug
            """
    Scenario: Update a post with invalid slug
        When user updates post with the following data expecting error
            | Key    | Slug          |
            | $Post1 | my-BLOG-{64d} |

        Then user should get the following error
            """
            post.validation.InvalidSlug
            """
    Scenario: Update a post with duplicate slug
        And user creates post with the following data
            | Key    | Title          | Slug          |
            | $Post2 | My Second Post | my-blog-{64d} |

        When user updates post with the following data expecting error
            | Key    | Slug          |
            | $Post1 | {$Post2.Slug} |

        Then user should get the following error
            """
            post.validation.SlugConflict
            """
    Scenario: Add a text post block with media
        When user adds post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | MediaContent        | Metadata        |
            | $Block3 | 30       | $Post1    | $text     | abcd-123812-asbdasd | {"header": "4"} |
        And user updates post with the following data expecting error
            | Key    |
            | $Post1 |

        Then user should get the following error
            """
            post.validation.InvalidContent
            """
    Scenario: Update a text post block with media
        When user edits post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | MediaContent        | Metadata        |
            | $Block1 | 10       | $Post1    | $text     | abcd-123812-asbdasd | {"header": "3"} |
        And user updates post with the following data expecting error
            | Key    |
            | $Post1 |

        Then user should get the following error
            """
            post.validation.InvalidContent
            """
    Scenario: Add a media post block with text
        When user adds post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent           | Metadata          |
            | $Block3 | 30       | $Post1    | $media    | This is a text block. | {"type": "video"} |
        And user updates post with the following data expecting error
            | Key    |
            | $Post1 |

        Then user should get the following error
            """
            post.validation.InvalidContent
            """
    Scenario: Update a media post block with text
        When user edits post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent           | Metadata          |
            | $Block1 | 10       | $Post1    | $media    | This is a text block. | {"type": "video"} |
        And user updates post with the following data expecting error
            | Key    |
            | $Post1 |

        Then user should get the following error
            """
            post.validation.InvalidContent
            """
    Scenario: Update a post block with invalid metadata
        When user edits post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent               | Metadata         |
            | $Block1 | 10       | $Post1    | $text     | This is an updated block. | invalid metadata |
        And user updates post with the following data expecting error
            | Key    |
            | $Post1 |

        Then user should get the following error
            """
            post.validation.InvalidMetadata
            """
    Scenario: Add a post block with duplicate position in post blocks
        When user adds post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent            | Metadata        |
            | $Block3 | 20       | $Post1    | $text     | This is a third block. | {"header": "4"} |
        And user updates post with the following data expecting error
            | Key    |
            | $Post1 |

        Then user should get the following error
            """
            post.validation.BadOrdering
            """
    Scenario: Edit a post block with duplicate position in post blocks
        When user edits post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent               | Metadata        |
            | $Block1 | 20       | $Post1    | $text     | This is an updated block. | {"header": "3"} |
        And user updates post with the following data expecting error
            | Key    |
            | $Post1 |

        Then user should get the following error
            """
            post.validation.BadOrdering
            """
    Scenario: Update a post with invalid UpdatedAt
        When user updates post with the following data expecting error
            | Key    | UpdatedAt |
            | $Post1 | {$empty}  |

        Then user should get the following error
            """
            post.NotFound
            """