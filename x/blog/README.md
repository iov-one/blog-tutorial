# Blog module

## Requirements

This module defines the required components for blog.

- A blog is where a user posts their article
- Every user can post article on their blog and has permission delete only their article
- Blog owner can set a time to delete the article during and after creation
- Users can like and comment to others articles

### State

- #### User

  - ID
  - Username
  - Bio

- #### Blog

  - ID
  - Owner
  - Title
  - Description
  - CreatedAt

- #### Article

  - ID
  - BlogID
  - Title
  - Content
  - CreatedAt
  - DeleteAt
  - CommentCount
  - LikeCount

- #### Comment
  
  - ID
  - ArticleID
  - Owner
  - Content

- #### Like
  
  - ID
  - ArticleID
  - Owner

### Messages

- #### Create User

  - Username
  - Bio

- #### Create Blog

  - Title
  - Description

- #### Create Article

  - BlogID
  - Title
  - Content
  - DeleteAt

- #### Create Comment
  
  - ArticleID
  - Content

- #### Like Article
  
  - ArticleID
