Feature: Post Create
    Scenario: Create a post
        Given user defines post with the following data
            | Key    | Title         | Slug          |
            | $Post1 | My First Post | my-blog-{64d} |
        And user adds post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent                 | Metadata        |
            | $Block1 | 10       | $Post1    | $text     | This is a text block.       | {"header": "3"} |
            | $Block2 | 20       | $Post1    | $text     | This is another text block. | {"header": "3"} |
        When user creates post with the following data
            | Key    |
            | $Post1 |

        Then user should be able to see post with the following data
            | Key    |
            | $Post1 |
    Scenario: Create a post with empty title
        Given user defines post with the following data
            | Key    | Title    | Slug          |
            | $Post1 | {$empty} | my-blog-{64d} |
        And user adds post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent           | Metadata        |
            | $Block1 | 10       | $Post1    | $text     | This is a text block. | {"header": "3"} |
        When user creates post with the following data expecting error
            | Key    |
            | $Post1 |

        Then user should get the following error
            """
            post.validation.EmptyTitle
            """
    Scenario: Create a post with empty slug
        Given user defines post with the following data
            | Key    | Title         | Slug     |
            | $Post1 | My First Post | {$empty} |
        And user adds post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent           | Metadata        |
            | $Block1 | 10       | $Post1    | $text     | This is a text block. | {"header": "3"} |
        When user creates post with the following data expecting error
            | Key    |
            | $Post1 |

        Then user should get the following error
            """
            post.validation.EmptySlug
            """
    Scenario: Create a post with invalid slug
        Given user defines post with the following data
            | Key    | Title         | Slug          |
            | $Post1 | My First Post | my-BLOG-{64d} |
        And user adds post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent           | Metadata        |
            | $Block1 | 10       | $Post1    | $text     | This is a text block. | {"header": "3"} |
        When user creates post with the following data expecting error
            | Key    |
            | $Post1 |

        Then user should get the following error
            """
            post.validation.InvalidSlug
            """
    Scenario: Create a post with duplicate slug
        Given user creates post with the following data
            | Key    | Title         | Slug          |
            | $Post1 | My First Post | my-blog-{64d} |

        When user creates post with the following data expecting error
            | Key    | Title         | Slug          |
            | $Post2 | My First Post | {$Post1.Slug} |

        Then user should get the following error
            """
            post.validation.SlugConflict
            """
    Scenario: Create a text post block with media
        When user defines post with the following data
            | Key    | Title         | Slug          |
            | $Post1 | My First Post | my-blog-{64d} |
        And user adds post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | MediaContent        | Metadata        |
            | $Block1 | 10       | $Post1    | $text     | abcd-123812-asbdasd | {"header": "3"} |
        And user creates post with the following data expecting error
            | Key    |
            | $Post1 |

        Then user should get the following error
            """
            post.validation.InvalidContent
            """
    Scenario: Create a media post block with text
        When user defines post with the following data
            | Key    | Title         | Slug          |
            | $Post1 | My First Post | my-blog-{64d} |
        And user adds post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent           | Metadata          |
            | $Block1 | 10       | $Post1    | $media    | This is a text block. | {"type": "video"} |
        And user creates post with the following data expecting error
            | Key    |
            | $Post1 |

        Then user should get the following error
            """
            post.validation.InvalidContent
            """
    Scenario: Create a post block with invalid metadata
        When user defines post with the following data
            | Key    | Title         | Slug          |
            | $Post1 | My First Post | my-blog-{64d} |
        And user adds post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent           | Metadata         |
            | $Block1 | 10       | $Post1    | $media    | This is a text block. | invalid metadata |
        And user creates post with the following data expecting error
            | Key    |
            | $Post1 |

        Then user should get the following error
            """
            post.validation.InvalidMetadata
            """
    Scenario: Create a post with duplicate position in post blocks
        Given user defines post with the following data
            | Key    | Title         | Slug          |
            | $Post1 | My First Post | my-blog-{64d} |
        And user adds post blocks with the following data
            | Key     | Position | HeaderKey | BlockType | TextContent                    | Metadata        |
            | $Block1 | 10       | $Post1    | $text     | This is the first text block.  | {"header": "3"} |
            | $Block2 | 10       | $Post1    | $text     | This is the second text block. | {"header": "3"} |
        When user creates post with the following data expecting error
            | Key    |
            | $Post1 |

        Then user should get the following error
            """
            post.validation.BadOrdering
            """